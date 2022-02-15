import http from 'k6/http';
import {check, randomSeed, sleep} from "k6";

export default function () {
    const url = "http://136.159.209.204:9099";
    const payload = JSON.stringify({
            args: [
                "--memrate",
                "3",
                "--memrate-bytes",
                "5G"
            ],
            timeout: "90s"
        })

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post(url + "/stress", payload, params);
    sleep(1000 * 10);
}

// import http from "k6/http";
// import {check, randomSeed, sleep} from "k6";
// import {Counter, Trend} from 'k6/metrics';
//
// export let options = {
//     vus: argsvus,
//     duration: 'argsduration',
//     userAgent: 'MyK6UserAgentString/1.0',
// };
//
// const BaseURL = "http://136.159.209.204:9099";
//
// let s1_duration_trend = new Trend('s1_duration_trend');
// // let edit_book_duration_trend = new Trend('edit_book_duration_trend');
// // let get_book_duration_trend = new Trend('get_book_duration_trend');
//
// let s1_counter = new Counter('s1_counter');
// // let edit_book_counter = new Counter('edit_book_counter');
// // let get_book_counter = new Counter('get_book_counter');
//
// export default function (data) {
//
//     let uniqueNumber = __VU * 100000000 + (__ITER % 4000);
//     randomSeed(uniqueNumber);
//
//     // const LoginAuthProb = argsloginprob;
//     // const GetBookProb = argsgetbookprob;
//     // const EditBookProb = argseditbookprob;
//     const SLEEP_DURATION = argssleepduration;
//
//     function get_random_item(items) {
//         const item = items[Math.floor(Math.random() * items.length)];
//         return item
//     }
//
//
//     const execute_random_stress = function () {
//
//         let body = JSON.stringify({
//             args: [
//                 "--memrate",
//                 "3",
//                 "--memrate-bytes",
//                 "5G"
//             ],
//             timeout: options.duration
//         });
//
//         let params = {
//             headers: {
//                 'Content-Type': 'application/json',
//                 'debug_id': new Date().getTime(),
//             }
//         };
//
//         let response = http.post(
//             BaseURL + '/stress',
//             body,
//             params
//         );
//         s1_duration_trend.add(response.timings.duration);
//         s1_counter.add(1);
//         check(response, {
//             'is_s1_200': r => r.status === 200,
//         });
//     }
//
//     // const r = Math.random();
//     // const sTime = Math.random() * SLEEP_DURATION + 0.5 * SLEEP_DURATION;
//     execute_random_stress();
//     sleep(SLEEP_DURATION);
//     // if (r < LoginAuthProb) {
//     // } else if (r >= LoginAuthProb && r < LoginAuthProb + GetBookProb) {
//     //     execute_get_book();
//     // } else {
//     //     execute_edit_book();
//     // }
//
// };
//
// export function teardown(data) {
// }
