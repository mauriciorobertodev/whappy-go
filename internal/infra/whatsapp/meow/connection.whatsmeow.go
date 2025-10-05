package meow

import (
	"context"
	"fmt"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (g *WhatsmeowGateway) Ping(ctx context.Context, inst *instance.Instance) (whatsapp.Ping, error) {
	if inst.Status.IsBanned() {
		return whatsapp.Ping{
			IsLoggedIn:  false,
			IsConnected: false,
		}, nil
	}

	if !inst.Status.IsLoggedIn() {
		return whatsapp.Ping{
			IsLoggedIn:  false,
			IsConnected: false,
		}, nil
	}

	client, ok := g.getClient(inst.ID)
	if ok {
		return whatsapp.Ping{
			IsLoggedIn:  client.IsLoggedIn(),
			IsConnected: client.IsConnected(),
		}, nil
	}

	client = whatsmeow.NewClient(g.container.NewDevice(), nil)

	if err := client.Connect(); err != nil {
		return whatsapp.Ping{
			IsLoggedIn:  false,
			IsConnected: false,
		}, nil
	}

	time.Sleep(2 * time.Second) // Wait a bit to ensure the connection is established

	ping := whatsapp.Ping{
		IsLoggedIn:  client.IsLoggedIn(),
		IsConnected: client.IsConnected(),
	}

	client.Disconnect()

	return ping, nil
}

func (g *WhatsmeowGateway) Connect(ctx context.Context, inst *instance.Instance) error {
	client, ok := g.getClient(inst.ID)
	if !ok {
		if inst.HasLoggedInBefore() {
			device, _ := types.ParseJID(inst.Device)

			deviceStore, err := g.container.GetDevice(ctx, device)
			if err != nil {
				return err
			}

			// If the device store exists, maybe the instance was already paired, so we try connect
			if deviceStore != nil {
				client := whatsmeow.NewClient(deviceStore, nil)

				client.AddEventHandler(func(evt interface{}) {
					switch v := evt.(type) {
					case *events.QR:
						fmt.Println("QR code received, please scan it")
					case *events.PairSuccess:
						fmt.Println("Pairing successful!")
					case *events.PairError:
						fmt.Println("Pairing error")
					case *events.QRScannedWithoutMultidevice:
						fmt.Println("QR code scanned without multi-device support!")
					case *events.Connected:
						fmt.Println("Connected to WhatsApp Web server!")
					case *events.KeepAliveTimeout:
						fmt.Println("Keep alive timeout!")
					case *events.KeepAliveRestored:
						fmt.Println("Keep alive restored!")
					case *events.LoggedOut:
						fmt.Println("Logged out from WhatsApp Web server!")
					case *events.StreamReplaced:
						fmt.Println("Stream was replaced, this usually means your phone logged in somewhere else!")
					case *events.TemporaryBan:
						fmt.Println("User got a temporary ban!", v.Code, v.Expire)
					case *events.HistorySync:
						fmt.Println("History sync completed!", v.Data.SyncType)
					case *events.UndecryptableMessage:
						fmt.Println("Received an undecryptable message!")
					case *events.NewsletterMessageMeta:
						fmt.Println("Received a newsletter message meta!")
					case *events.Message:
						fmt.Println("New message received from:", v.Info.Sender.String(), "in chat:", v.Info.Chat.String(), "with id:", v.Info.ID)

						if IsStatus(v.Info.Chat) {
							fmt.Println("Status message received")
							g.emitStatusNew(inst, v)
							return
						}

						if IsGroup(v.Info.Chat) {
							groupType := g.getGroupType(ctx, inst, v.Info.Chat.String())
							if groupType == group.GroupTypeAnnouncement {
								fmt.Println("Announcement group message received")
								g.emitCommunityAnnouncement(inst, v)
								return
							}

							fmt.Println("Group message received")
							g.emitGroupNewMessage(inst, v)
							return
						}

						if IsNewsletter(v.Info.Chat) {
							fmt.Println("Newsletter message received")
							g.emitNewsletterNewPost(inst, v)
							return
						}

						if IsUser(v.Info.Chat) {
							fmt.Println("User message received")
							g.emitUserNewMessage(inst, v)

							go (func() {
								time.Sleep(2 * time.Second) // Wait a bit to ensure the message is fully processed

								client.SendMessage(ctx, v.Info.Chat, &waE2E.Message{
									ReactionMessage: &waE2E.ReactionMessage{
										Text: proto.String("👍"),
										Key: &waCommon.MessageKey{
											ID:        proto.String(v.Info.ID),
											RemoteJID: proto.String(v.Info.Chat.String()),
											FromMe:    proto.Bool(false),
										},
									},
								})
							})()

							return
						}

						fmt.Println("Unknown message received")
						// TODO: maybe onUnknownMessage ?
					case *events.Receipt:
						switch v.Type {
						case types.ReceiptTypeDelivered:
							g.emitMessageDelivered(inst, v)
						case types.ReceiptTypeRead:
							g.emitMessageRead(inst, v)
						case types.ReceiptTypePlayed:
							g.emitMessagePlayed(inst, v)
						case types.ReceiptTypeRetry:
						case types.ReceiptTypeServerError:
						default:
						}
					case *events.ChatPresence:
						g.emitChatPresence(inst, v)
					case *events.Presence:
						g.emitUserPresence(inst, v)
					case *events.Disconnected:
						fmt.Println("Disconnected from WhatsApp Web server!")
					case *events.AppStateSyncComplete:
						fmt.Println("App state sync completed!")
					case *events.JoinedGroup:
						fmt.Println("Joined a group!")
					case *events.GroupInfo:
						g.emitGroupChangedInfo(inst, v)
					case *events.Picture:
						if IsGroup(v.JID) {
							groupType := g.getGroupType(ctx, inst, v.JID.String())
							if groupType == group.GroupTypeCommunity {
								g.emitCommunityChangedPhoto(inst, v)
								return
							}

							g.emitGroupChangedPhoto(inst, v)
							return
						}

						if IsNewsletter(v.JID) {
							g.emitNewsletterChangedPhoto(inst, v)
							return
						}

						if IsUser(v.JID) {
							g.emitUserChangedPhoto(inst, v)
							return
						}

						// TODO: maybe onUnknownPhoto ?
					case *events.DeleteChat:
						g.emitChatDeleted(inst, v)
					case *events.ClearChat:
						g.emitChatCleared(inst, v)
					case *events.MarkChatAsRead:
						g.emitChatRead(inst, v)
					case *events.Pin:
						g.emitChatChangedPin(inst, v)
					case *events.Mute:
						g.emitChatChangedMute(inst, v)
					case *events.Archive:
						g.emitChatChangedArchive(inst, v)
					case *events.UserAbout:
						g.emitUserChangedStatus(inst, v)
					case *events.PushNameSetting:
						fmt.Println("User push name changed!")
					case *events.IdentityChange:
						fmt.Println("User identity changed!")
					case *events.PrivacySettings:
						g.emitPrivacyChanged(inst, v)
					case *events.NewsletterJoin:
						fmt.Println("User joined a newsletter!")
					case *events.NewsletterLeave:
						fmt.Println("User left a newsletter!")
					case *events.NewsletterMuteChange:
						fmt.Println("User changed newsletter mute!")
					case *events.NewsletterLiveUpdate:
						fmt.Println("Newsletter live update!")
					case *events.Blocklist:
						g.emitBlocklistChanged(inst, v)
					case *events.Contact:
						fmt.Println("User contact updated!")
					default:
						fmt.Printf("Unknown event type %T\n", v)
					}
				})

				// If the device store exists, maybe the instance was already paired, so we try connect
				if err := client.Connect(); err != nil {
					return err
				}

				time.Sleep(2 * time.Second) // Wait a bit to ensure the connection is established

				if client.IsLoggedIn() && client.IsConnected() {
					g.clients.Store(inst.ID, client)
					return nil
				}
			}

			return fmt.Errorf("device store not found for instance that has logged in before, maybe the database was cleaned")
		}

		return fmt.Errorf("the instance has never logged in, it needs to login first")
	}

	if err := client.Connect(); err != nil {
		return err
	}

	return nil
}

func (g *WhatsmeowGateway) Disconnect(ctx context.Context, inst *instance.Instance) error {
	client, ok := g.getClient(inst.ID)
	if !ok {
		return instance.ErrInstanceNotConnected
	}

	if client.IsConnected() {
		client.Disconnect()
	}

	return nil
}

func (g *WhatsmeowGateway) getGroupType(ctx context.Context, inst *instance.Instance, jid string) group.GroupType {
	cacheKey := cache.CacheKeyGroupTypePrefix + jid

	if data, err := g.cache.Get(cacheKey); err == nil {
		return group.GroupType(string(data))
	}

	gp, err := g.GetGroup(ctx, inst, jid, utils.BoolPtr(false))
	if err != nil {
		return group.GroupTypeRegular
	}

	g.cache.Forever(cacheKey, []byte(gp.Type))

	return gp.Type
}
