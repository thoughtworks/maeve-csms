import ws from 'k6/ws';
import { check } from 'k6';
export const options = {
    duration: '10s',
    vus: 1
};

export default function () {
    const url = 'ws://localhost/ws/cs001';
    const params = {
        headers: {'Sec-WebSocket-Protocol': 'ocpp1.6', 'Authorization': 'Basic Y3MwMDE6ZmlkZGxlc3RpY2tzX2Zpc2hzdGlja3M='}
    };


    const res = ws.connect(url, params, function (socket) {
        socket.on('open', () => console.log('connected'));
        socket.on('message', () => console.log('Message received: ', '[2,"9","BootNotification",{"chargePointModel":"me100","chargePointVendor":"me"}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "10", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Available"}]'));
        socket.on('message', () => console.log('Message received: ', '[2,"11","Heartbeat", {}]'));
        socket.on('message', () => console.log('Message received: ', '[2,"12","Heartbeat", {}]'));
        socket.on('message', () => console.log('Message received: ', '[2,"13","Heartbeat", {}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "14", "[2, "16", "Authorize", {"idTag": "38748383L7337848H823"}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "15", "StartTransaction",{"connectorId": 1, "idTag": "38748383L7337848H823", "meterStart": 3, "reservationId": 5, "timestamp":"2023-09-18T08:25:40.20Z"}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "16", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "17", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]'));
        socket.on('message', () => console.log('Message received: ', '[2, "18", "StopTransaction", {"idTag": "38748383L7337848H823", "meterStop": 3, "timestamp": "2023-09-18T08:25:40.20Z", "transactionId": 3}]'));
        socket.on('close', () => console.log('disconnected'));
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}