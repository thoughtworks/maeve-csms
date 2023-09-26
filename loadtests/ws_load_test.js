import ws from 'k6/ws';
import { check } from 'k6';
export const options = {
    vus: 1,
    duration: '10s',
};

const csIds = [ 1, 2, 3, 4, 5];

export default function () {
    let csId = csIds[Math.floor(Math.random()*csIds.length)];
    const params = {
        headers: {'Sec-WebSocket-Protocol': 'ocpp1.6', 'Authorization': 'Basic Y3MwMDE6ZmlkZGxlc3RpY2tzX2Zpc2hzdGlja3M='}
    };

        const url = `ws://localhost/ws/cs00${csId}`;
        const res = ws.connect(url, params, function (socket) {
        socket.send('[2,"1","BootNotification",{"chargePointModel":"me100","chargePointVendor":"me"}]');
        socket.send('[2,"2","BootNotification",{"chargePointModel":"me100","chargePointVendor":"me"}]');
        socket.send('[2, "3", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Available"}]');
        socket.send('[2,"4","Heartbeat", {}]');
        socket.send('[2,"5","Heartbeat", {}]');
        socket.send('[2,"6","Heartbeat", {}]');
        socket.send('[2, "7", "[2, "16", "Authorize", {"idTag": "38748383L7337848H823"}]');
        socket.send('[2, "8", "StartTransaction",{"connectorId": 1, "idTag": "38748383L7337848H823", "meterStart": 3, "reservationId": 5, "timestamp":"2023-09-18T08:25:40.20Z"}]');
        socket.send('[2, "9", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]');
        socket.send('[2, "10", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]');
        socket.send('[2, "11", "StopTransaction", {"idTag": "38748383L7337848H823", "meterStop": 3, "timestamp": "2023-09-18T08:25:40.20Z", "transactionId": 3}]');
        socket.close()
    }
    );

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}