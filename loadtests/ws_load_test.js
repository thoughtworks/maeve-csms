import ws from 'k6/ws';
import { check } from 'k6';
import encoding from 'k6/encoding';
import exec from 'k6/execution';
import { sleep } from 'k6';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'ramping-vus',
            startVUs: 2,
            stages: [
                { duration: '5m', target: 50 },
                { duration: '5m', target: 100 },
                { duration: '10m', target: 200 },
                { duration: '10m', target: 300 },
                { duration: '10m', target: 400 },
                { duration: '10m', target: 500 },
                { duration: '10m', target: 600 },
                { duration: '5m', target: 550 },
                { duration: '5m', target: 600 },
                { duration: '5m', target: 550 },
                { duration: '5m', target: 600 },
                { duration: '10m', target: 500 },
                { duration: '10m', target: 400 },
                { duration: '10m', target: 300 },
                { duration: '10m', target: 100 },
                { duration: '5m', target: 50 },
                { duration: '5m', target: 0 }
            ],
            gracefulRampDown: '0s',
        },
    },
};

export default function () {
    let vuIdInTest = exec.vu.idInTest
    let data = `cs${vuIdInTest}:fiddlesticks_fishsticks`;
    let base64data = encoding.b64encode(data)
    const params = {
        headers: {'Sec-WebSocket-Protocol': 'ocpp1.6', 'Authorization': `Basic ${base64data}`}
    };
        const url = `ws://localhost/ws/cs${vuIdInTest}`;
        const res = ws.connect(url, params, function (socket) {
        socket.send('[2,"1","BootNotification",{"chargePointModel":"me100","chargePointVendor":"me"}]');
        sleep(10)
        socket.send('[2,"2", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Available"}]');
        sleep(10)
        socket.send('[2,"3","Heartbeat", {}]');
        sleep(10)
        socket.send('[2,"4","Heartbeat", {}]');
        sleep(10)
        socket.send('[2,"5","Heartbeat", {}]');
        sleep(10)
        socket.send('[2, "6", "Authorize", {"idTag": "38748383L7337848H823"}]');
        sleep(10)
        socket.send('[2,"7", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Preparing"}]');
        sleep(10)
        socket.send('[2,"8", "StartTransaction",{"connectorId": 1, "idTag": "38748383L7337848H823", "meterStart": 3, "reservationId": 5, "timestamp":"2023-09-18T08:25:40.20Z"}]');
        sleep(300)
        socket.send('[2,"9","Heartbeat", {}]');
        sleep(60)
        socket.send('[2,"10","Heartbeat", {}]');
        sleep(60)
        socket.send('[2,"11","Heartbeat", {}]');
        sleep(60)
        sleep(300)
        socket.send('[2,"12", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]');
        sleep(60)
        socket.send('[2,"13", "MeterValues", {"connectorId": 1, "meterValue":[{"timestamp":"2023-09-18T08:25:40.20Z", "sampledValue": [{"value": "5"}]}]}]');
        sleep(60)
        socket.send('[2,"14", "StopTransaction", {"idTag": "38748383L7337848H823", "meterStop": 3, "timestamp": "2023-09-18T08:25:40.20Z", "transactionId": 3}]');
        sleep(10)
        socket.send('[2,"15", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Finishing"}]');
        sleep(10)
        socket.send('[2,"16", "StatusNotification", {"connectorId": 1, "errorCode": "NoError", "status": "Available"}]');
        sleep(10)
        socket.close()
    }
    );

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}