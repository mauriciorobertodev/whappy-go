# Whappy GO

## WhatsApp + Happy + Go 👍  
A simple WhatsApp HTTP API built on top of **whatsmeow**.

## 🎯 What is Whappy GO?  
**Whappy GO** is a lightweight HTTP API that wraps **whatsmeow**, providing a clean and consistent way to interact with WhatsApp.  
It makes it easy to **send messages**, **manage instances**, and **handle groups** through RESTful endpoints.

## 🧑‍💻 About the Author  
Hey there! I’m **Mauricio Roberto**, a Junior Fullstack Developer from **Brazil 🇧🇷**.  
I built **Whappy GO** as a side project to learn more about **Golang** (I usually work with **Laravel** and **TypeScript**), explore **DDD** concepts (though I’ll admit, that part didn’t go so well 😅),  
and, of course, to level up my **English** — which is still pretty tough for me 😩.
<br/>

## ⚡ Features  

- 🚀 **HTTP Endpoints for WhatsApp** — simple and intuitive REST API for message automation.  
- 🔐 **Multi-instance Authentication** — manage multiple WhatsApp sessions via QR Code.  
- 💬 **Messaging Support** — send text, images, documents, audio, and more.  
- 👥 **Group & Contact Management** — create, update, and manage groups and contacts.  
- ⚡ **High-performance Core** — built with [Fiber](https://gofiber.io/) for fast and efficient HTTP handling.  
- 🗄️ **Database Drivers** — support for `sqlite` and `postgres` with seamless integration.  
- 📦 **Storage System** — choose between `local` or `S3` for global automatic media saving, and uploads feature.  
- 📤 **Uploads System** — enable upload routes for users, store files once, and reference them by upload ID.
- 🕋 **Cache Layer** — `in-memory` or `redis` caching for fast lookups and reduced load.  
- 📤 **Upload Cache** — configurable cache for WhatsApp server uploads (default: 24 h).  
- 🧩 **Flexible Authentication** — use instance tokens or impersonate an instance via `ADMIN_TOKEN` + `X-Instance-ID` header.
- 📝 **Beautiful Documentation** — clear API reference and a polished web interface 😏.
- 🛠 **Event Bus System** — central event hub with `in-memory` and `redis` Pub/Sub drivers for flexible events consumption.
<br/>

## 📌 Endpoints

View the full [API Documentation](https://go.whappy.com.br) — proudly made with [API Dog](https://app.apidog.com/invite/user?token=cjIKD_sGAqQUMrt6-F9KX).

### Legend
✅ **Done** | 🚧 **In Progress** | ❌ **Not Started**

---

### 📦 Instances

> **Note:** This endpoints need `Authorization` header with `ADMIN_TOKEN`.

- ✅ **GET**    `/admin/instances`            – List all instances.  
- ✅ **POST**   `/admin/instances`            – Create a new instance. 
- ✅ **GET**    `/admin/instances/{id}`       – Get instance details.   
- ❌ **PUT**    `/admin/instances/{id}`       – Update an instance.   
- ❌ **DELETE** `/admin/instances/{id}`       – Delete an instance.    
- ✅ **PUT**    `/admin/instances/{id}/token` – Renew a token instance.  

> **Note:** All endpoints above here require the `Authorization` header with the instance token. OR you can send `Authorization` with `ADMIN_TOKEN` + `X-INSTANCE-ID` with the instance ID.

### 🔐 Auth
Endpoints to pair/unpair an whatsapp account.

- ✅ **POST** `/session/login`  – Log out the instance (requires authentication to send messages).   
- ✅ **POST** `/session/logout` – Log out the instance (requires authentication to send messages).   
- ✅ **GET**  `/session/qr`     – Generate QR code to connect WhatsApp.   

### 🔌 Connection

Endpoints to connect, disconnect, or check the instance via whatsApp webSocket.

✅ **POST** `/session/connect`    – Connect the WhatsApp server instance (events are now being listened to).   
✅ **POST** `/session/disconnect` – Disconnect the WhatsApp server instance.   
✅ **GET**  `/session/ping`       – 

### 🧑‍💻 Users
❌ **GET** `/users/{jid or lid}/info`.    –   
❌ **GET** `/users/{jid or lid}/photo`    –  
❌ **GET** `/users/{jid or lid}/presence` –  

### 📨 Messages

Endpoints to send messages.

✅ **GET** `/messages/id` – Generate message IDs whatsapp like, multi id can be generated using `?quantity=8`.  

✅ **POST** `/messages/text`     – Send text message.  
✅ **POST** `/messages/image`    – Send image message.  
✅ **POST** `/messages/video`    – Send video message.  
✅ **POST** `/messages/audio`    – Send audio message.  
✅ **POST** `/messages/voice`    – Send voice message.  
❌ **POST** `/messages/sticker`  – Send sticker message.  
❌ **POST** `/messages/location` – Send location message.  
❌ **POST** `/messages/contact`  – Send contact message.  
❌ **POST** `/messages/gif`      – Send gif message.  
❌ **POST** `/messages/poll`     – Send poll message.  
✅ **POST** `/messages/reaction` –   

✅ **POST** `/messages/read` – Mark messages as read. (many messages supported).  


### 👤 Contacts

Endpoints to manage contacts.

✅ **GET**  `/contacts`                – List all contacts.    
✅ **GET**  `/contacts/{phone or jid}` – Get details of a contact.     
✅ **POST** `/contacts/check`          – Check if given phone numbers exist on WhatsApp.  

### 🚫 Blocklist

Endpoints to manage the blocklist.

✅ **GET**    `/blocklist`                –   
✅ **POST**   `/blocklist/{phone or jid}` –    
✅ **DELETE** `/blocklist/{phone or jid}` –  

### 👥 Groups

Endpoints to manage groups.

✅ **GET**    `/groups`      – List joined groups and participants.   
✅ **POST**   `/groups`      – Create a new group.  
✅ **GET**    `/groups/{id}` – Get group info.   
✅ **PATCH**  `/groups/{id}` – Update group permissions.    
✅ **DELETE** `/groups/{id}` – Leave group.  

✅ **PATCH** `/groups/{id}/name`         – Update group name.  
✅ **PATCH** `/groups/{id}/description`  – Update group description. 
✅ **PATCH** `/groups/{id}/disappearing` – Update message disappearing settings. 

✅ **GET**    `/groups/{id}/photo` – Get group photo.  
✅ **PUT**    `/groups/{id}/photo` – Update group photo.    
✅ **DELETE** `/groups/{id}/photo` – Delete group photo.  

✅ **POST** `/groups/join` – Enter on group.  

✅ **GET**    `/groups/{id}/invite` – Get group invite link.  
✅ **DELETE** `/groups/{id}/invite` – Revoke group invite link and return new link. 

✅ **GET**    `/groups/{id}/participants` – Get participants.    
✅ **POST**   `/groups/{id}/participants` – Add participants.   
✅ **DELETE** `/groups/{id}/participants` – Remove participants.   

✅ **POST**   `/groups/{id}/admins` – Promote participants to admin.   
✅ **DELETE** `/groups/{id}/admins` – Demote admins.  

### 👥 Communities

Endpoints to manage communities.

❌ **GET**    `/communities`             – List joined communities and groups with participants.    
❌ **POST**   `/communities`             – Create a new community.  
❌ **GET**    `/communities/{id}`        – Get community info.    
❌ **POST**   `/communities/{id}/groups` – Links groups.  
❌ **DELETE** `/communities/{id}/groups` – Unlinks groups.   

### 💬 Chat

Endpoints with utils for chats.

✅ **POST**  `/chat/presence` – Change presence in chat to TYPING/RECORDING/PAUSE.   
❌ **PATCH** `/chat/mute`     –   
❌ **PATCH** `/chat/pin`      –   

### 📸 Pictures

Endpoints to fetch pictures.

✅ **GET** `/pictures/{phone or jid}` –   

### 📤 Uploads

**Works only if storage is configured.**
Endpoints to manage uploads, used later when sending messages.

✅ **GET**    `/uploads`      – List stored files.    
✅ **POST**   `/uploads`      – Upload.  
✅ **PUT**    `/uploads/{id}` – Update.  
✅ **GET**    `/uploads/{id}` – Get.  
✅ **DELETE** `/uploads/{id}` – Delete. 

### 💅 Status

Endpoints to manage status.

❌ **POST** `/status/text`  – Create a text status.    
❌ **POST** `/status/image` – Create a image status.  
❌ **POST** `/status/audio` – Create an audio status.  
❌ **POST** `/status/video` – Create a video status.  

### ⬇️ Download
❌ **GET** `/download/image`    –   
❌ **GET** `/download/video`    –   
❌ **GET** `/download/audio`    –   
❌ **GET** `/download/sticker`  –   
❌ **GET** `/download/document` –   

### 🌐 Webhooks
✅ **GET**    `/webhooks`      – Get all webhooks.  
✅ **POST**   `/webhooks`      – Create a new webhook.  
✅ **GET**    `/webhooks/{id}` – Get a specific webhook.  
✅ **PUT**    `/webhooks/{id}` – Update a specific webhook.  
✅ **DELETE** `/webhooks/{id}` – Delete a specific webhook.  

<br/>

## 💻 API Clients / SDKs
Here you can list SDKs, libraries, or clients that integrate with Whappy GO.

🚧 **Whappy GO Laravel SDK** – A PHP client built using Saloon for easy integration with Whappy GO. By @mauriciorobertodev 

✨ *(Add your own SDK or client here)* – e.g., Node.js, Python, Golang, etc.
<br/>

## 🛠️ Built With Whappy GO
Show off your projects, tools, or services built with Whappy GO.
Feel free to add anything you’ve created using the API.

<!-- 💬 **Whappy Desktop** – Manage contacts, organize lists, and plan message campaigns with ease. @mauriciorobertodev -->

✨ *(Add your project here)* – Share what you've built with Whappy GO.
<br/>

## 💬 Join the Community

Have questions, ideas, or feedback?  
Come chat with us in the [Discussions](https://github.com/mauriciorobertodev/whappy-go/discussions)!

You can also:
- 🐛 Report bugs in [Issues](https://github.com/mauriciorobertodev/whappy-go/issues)
- 💡 Suggest features in the [Ideas](https://github.com/mauriciorobertodev/whappy-go/discussions/categories/ideas)
- Help improve docs or examples ✨
<br/>

## 👩‍💻 How to Contribute

Contributions are always welcome!  
If you don’t have a specific idea in mind, you can look for `TODO:` comments in the code,  
or help improve the **documentation**, fix **bugs**, or add **new MIME type extensions**.

1. **Fork** the repository  
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)  
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)  
4. **Push** to your branch (`git push origin feature/amazing-feature`)  
5. **Open** a Pull Request
<br/>

## 🤝 Contributors
<a href="https://github.com/mauriciorobertodev/whappy-go/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=mauriciorobertodev/whappy-go" />
</a>
  
<br/>

## 🌟 Support

If you find **Whappy GO** useful, please consider:
- ⭐ Starring the repository
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 📖 Improving [API Docs](https://go.whappy.com.br)
- 🧑‍💻 Contributing code
- 🧪 Create tests
<br/>

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/mauriciorobertodev/whappy-go/blob/main/LICENSE.md) file for details.
<br/>


## 🔬 References

While building **Whappy GO**, I studied several resources that helped me structure ideas around architecture, usability, and design patterns.

- [Wuzapi inspiration](https://github.com/asternic/wuzapi) [@asternic/wuzapi](https://github.com/asternic/wuzapi)
- [Whatsmeow documentation](https://pkg.go.dev/go.mau.fi/whatsmeow#example-package) [@tulir/whatsmeow](https://github.com/tulir/whatsmeow)
- [Fiber documentation](https://docs.gofiber.io/next/) [@gofiber/fiber](https://github.com/gofiber/fiber)
- [SQLX Documentation](https://jmoiron.github.io/sqlx/) [@jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- [Ginkgo documentation](https://onsi.github.io/ginkgo/#getting-started) [@onsi/ginkgo](https://github.com/onsi/ginkgo)
- [Webhook Cool](https://webhook.cool/at/microscopic-mother-55/zIZuen9CfRwIYGUefqux4tKxf9WXn8Oj)
- [Video: Authorization: Domain or Application Layer?](https://www.youtube.com/watch?v=0TpejBzN-xw)  
- [Video: Domain Driven Design (DDD) in Golang!?](https://www.youtube.com/watch?v=6FY9urgIjqo)  
- [Video: cursor-based vs offset-based pagination](https://www.youtube.com/watch?v=a3GgMVeCwHg)  
- [Article: Golang Functional Options Pattern](https://golang.cafe/blog/golang-functional-options-pattern.html)
- [Article: Webhook security](https://hookdeck.com/webhooks/guides/complete-guide-to-webhook-security)
- [Chat GPT](https://chatgpt.com/)
