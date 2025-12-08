const WebSocket = require('ws');

class WSMockClient {
  constructor(url, clientId) {
    this.clientId = clientId;
    this.ws = new WebSocket(url);
    
    this.ws.on('open', () => {
      console.log(`Cliente ${clientId} conectado`);
      this.ws.send(JSON.stringify({
        type: 'register',
        clientId: clientId
      }));
    });
    
    this.ws.on('message', (data) => {
      console.log(`Cliente ${clientId} recebeu:`, data.toString());
    });
    
    this.ws.on('close', () => {
      console.log(`Cliente ${clientId} desconectado`);
    });
  }
}

const clients = [];
for (let i = 0; i < 5; i++) {
  clients.push(new WSMockClient('ws://localhost:8080/echo', `client-${i}`));
}