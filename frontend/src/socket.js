const MESSAGE_ADD_CONTAINER = 1;
const MESSAGE_REMOVE_CONTAINER = 2;
const MESSAGE_FETCH_CONTAINERS = 3;

export default class Socket {
  constructor(url) {
    this.url = url;
    this.ws = null;
    this.onAddContainer = (c) => {};
    this.onRemoveContainer = (id) => {};
    this.onFetchContainers = (containers) => {};
  }

  open() {
    this.ws = new WebSocket(this.url);
    this.ws.onmessage = (event) => {
      const messages = event.data.split('\n');
      for (const message of messages) {
        const msg = JSON.parse(message);
        switch (msg.kind) {
          case MESSAGE_ADD_CONTAINER:
            this.onAddContainer(msg.payload);
          case MESSAGE_REMOVE_CONTAINER:
            this.onRemoveContainer(msg.payload);
          case MESSAGE_FETCH_CONTAINERS:
            this.onFetchContainers(msg.payload);
        }
      }
    };
  }

  onMessage(message) {
    console.log(message);
  }
}
