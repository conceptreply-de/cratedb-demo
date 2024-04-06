import http from 'k6/http';
import { sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 10,
  // A string specifying the total duration of the test run.
  duration: '600s',
};

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function() {
  http.post('http://localhost:8080/vehicles/1/measurements', JSON.stringify(
    {"value": randomIntBetween(43, 85), "type": "STATE_OF_CHARGE"}
  ), {
    headers: { 'Content-Type': 'application/json' },
  });
  sleep(1);
}
