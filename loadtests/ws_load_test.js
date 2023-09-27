import ws from 'k6/ws';
import { check } from 'k6';
import encoding from 'k6/encoding';
import exec from 'k6/execution';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '20s', target: 5 },
                { duration: '10s', target: 0 },
            ],
            gracefulRampDown: '0s',
        },
    },
};

export default function () {
    let vuIdInTest = exec.vu.idInTest
    let data = `cs00${vuIdInTest}:fiddlesticks_fishsticks`;
    let base64data = encoding.b64encode(data)
    const params = {
        headers: {'Sec-WebSocket-Protocol': 'ocpp1.6', 'Authorization': `Basic ${base64data}`}
    };
        const url = `ws://localhost/ws/cs00${vuIdInTest}`;
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