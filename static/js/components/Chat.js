import { App } from "../index.js";
import { $, makeElements, formatTimestamp, isoDate } from "../utils.js";
import { ws } from "../JS-handlers.js";

export class Chat {
  constructor(partnerId, partnerUsername) {
    this.partnerId = partnerId;
    this.partnerUsername = partnerUsername;
    this.sidebarElement = $("#sidebar");
    this.visible = false;
  }

  render() {
    const chatCard = makeElements({
      type: "div",
      name: `chatArea-${this.partnerId}`,
      classNames: ["chatArea", "hidden"],
      children: [
        makeElements({
          type: "div",
          classNames: "chat-header",
          children: [
            makeElements({
              type: "h3",
              contents: "Username",
              classNames: "username",
            }),
            makeElements({ type: "div", classNames: "close-chat" }),
          ],
        }),
        makeElements({ type: "div", classNames: "chat" }),
        makeElements({
          type: "div",
          classNames: "respond-container",
          children: [
            makeElements({
              type: "textarea",
              classNames: "chat-input",
              placeholder: "Write a message...",
              attributes: { resize: "none", spellcheck: "false" },
            }),
            makeElements({
              type: "button",
              classNames: "sendButton",
              contents: "Send",
            }),
          ],
        }),
      ],
    });

    this.sidebarElement.appendChild(chatCard);

    $(`#chatArea-${this.partnerId} .sendButton`).addEventListener(
      "click",
      () => {
        this.handleSendMessage();
      }
    );

    $(`#chatArea-${this.partnerId} .close-chat`).addEventListener(
      "click",
      () => {
        this.hide();
      }
    );
  }
  async loadConversation() {
    $(`#chatArea-${this.partnerId} .chat`).innerHTML = "";
    $(`#chatArea-${this.partnerId} .chat-header .username`).textContent =
      this.partnerUsername;

    ws.send(
      JSON.stringify({
        type: "get_conversations",
        recipient: String(this.partnerId),
        // test time stamp. replace with actual timestamp value of oldest message
        created_at: isoDate("17.10.2024 23:04:55"),
      })
    );
  }
  appendConversation(messages) {
    if (Array.isArray(messages)) {
      messages.forEach((message) => {
        const messageElement = makeElements({
          type: "div",
          classNames:
            message.sender === this.partnerId
              ? "message-received"
              : "message-sent",
          contents: `${message.sender_username}: ${message.content}`,
        });
        $(`#chatArea-${this.partnerId} .chat`).appendChild(messageElement);
      });
    } else {
      let senderUsername = messages.sender_username;
      if (!senderUsername) {
        if (messages.recipient == App.user.id) {
          senderUsername = this.partnerUsername;
        } else {
          senderUsername = App.user.username;
        }
      }
      const messageElement = makeElements({
        type: "div",
        classNames:
          messages.sender == App.user.id
            ? ["message", "message-sent"]
            : ["message", "message-received"],
        children: [
          makeElements({
            type: "div",
            classNames: "heading",
            children: [
              makeElements({
                type: "div",
                classNames: "message-author",
                contents: `${senderUsername}`,
              }),
              makeElements({
                type: "div",
                classNames: "message-time",
                contents: formatTimestamp(messages.created_at),
              }),
            ],
          }),
          makeElements({
            type: "div",
            classNames: "message-content",
            contents: `${messages.content}`,
          }),
        ],
      });
      $(`#chatArea-${this.partnerId} .chat`).appendChild(messageElement);
    }
  }
  show() {
    $(`#chatArea-${this.partnerId}`).classList.remove("hidden");
    this.visible = true;
  }
  hide() {
    $(`#chatArea-${this.partnerId}`).classList.add("hidden");
    this.visible = false;
  }
  handleSendMessage() {
    const msgInput = $(`#chatArea-${this.partnerId} .chat-input`);
    if (msgInput.value.trim() === "") return;

    console.log(
      `${App.user.id} sends "user-message" to user ${this.partnerId}: ${msgInput.value}`
    );

    const msg = JSON.stringify({
      type: "user-message",
      recipient: `${this.partnerId}`,
      content: `${msgInput.value}`,
    });
    ws.send(msg);
    msgInput.value = "";
  }
}
