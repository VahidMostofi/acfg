import http from "k6/http";
import {check, randomSeed, sleep} from "k6";
import {Counter, Trend} from 'k6/metrics';

export let options = {
    vus: argsvus,
    duration: 'argsduration',
    userAgent: 'MyK6UserAgentString/1.0',
};

//const BaseURL = "http://172.21.0.2:9080";
// const BaseURL = "http://172.21.0.3:9099";
// const BaseURL = "http://k8s-bookstor-ingressb-c30fd4ec3a-643144735.us-west-2.elb.amazonaws.com"
// const BaseURL = "http://172.27.0.2:9099";
// const BaseURL = "http://136.159.209.204:9099";
// const BaseURL = "http://172.24.0.2:9099";

const BaseURL = "http://136.159.209.204:9099";

let login_duration_trend = new Trend('login_duration_trend');
let edit_book_duration_trend = new Trend('edit_book_duration_trend');
let get_book_duration_trend = new Trend('get_book_duration_trend');

let login_counter = new Counter('login_counter');
let edit_book_counter = new Counter('edit_book_counter');
let get_book_counter = new Counter('get_book_counter');

export default function (data) {
    // TODO should these defs be in a VU?
    const books = [{"_id": "5e5218426a4cea0021cdbf9a"}, {"_id": "5e5218426a4cea0021cdbf9b"}, {"_id": "5e5218426a4cea0021cdbf9c"}, {"_id": "5e5218436a4cea0021cdbf9d"}, {"_id": "5e5218436a4cea0021cdbf9e"}, {"_id": "5e5218436a4cea0021cdbf9f"}, {"_id": "5e5218436a4cea0021cdbfa0"}, {"_id": "5e5218436a4cea0021cdbfa1"}, {"_id": "5e5218446a4cea0021cdbfa2"}, {"_id": "5e5218446a4cea0021cdbfa3"}, {"_id": "5e5218446a4cea0021cdbfa4"}, {"_id": "5e5218446a4cea0021cdbfa5"}, {"_id": "5e5218446a4cea0021cdbfa6"}, {"_id": "5e5218446a4cea0021cdbfa7"}, {"_id": "5e5218446a4cea0021cdbfa8"}, {"_id": "5e5218446a4cea0021cdbfa9"}, {"_id": "5e5218446a4cea0021cdbfaa"}, {"_id": "5e5218446a4cea0021cdbfab"}, {"_id": "5e5218446a4cea0021cdbfac"}, {"_id": "5e5218446a4cea0021cdbfad"}, {"_id": "5e5218446a4cea0021cdbfae"}, {"_id": "5e5218446a4cea0021cdbfaf"}, {"_id": "5e5218446a4cea0021cdbfb0"}, {"_id": "5e5218446a4cea0021cdbfb1"}, {"_id": "5e5218456a4cea0021cdbfb2"}, {"_id": "5e5218456a4cea0021cdbfb3"}, {"_id": "5e5218456a4cea0021cdbfb4"}, {"_id": "5e5218456a4cea0021cdbfb5"}, {"_id": "5e5218456a4cea0021cdbfb6"}, {"_id": "5e5218456a4cea0021cdbfb7"}, {"_id": "5e5218456a4cea0021cdbfb8"}, {"_id": "5e5218456a4cea0021cdbfb9"}, {"_id": "5e5218456a4cea0021cdbfba"}, {"_id": "5e5218456a4cea0021cdbfbb"}, {"_id": "5e5218456a4cea0021cdbfbc"}, {"_id": "5e5218456a4cea0021cdbfbd"}, {"_id": "5e5218456a4cea0021cdbfbe"}, {"_id": "5e5218456a4cea0021cdbfbf"}, {"_id": "5e5218456a4cea0021cdbfc0"}, {"_id": "5e5218456a4cea0021cdbfc1"}, {"_id": "5e5218456a4cea0021cdbfc2"}, {"_id": "5e5218466a4cea0021cdbfc3"}, {"_id": "5e5218466a4cea0021cdbfc4"}, {"_id": "5e5218466a4cea0021cdbfc5"}, {"_id": "5e5218466a4cea0021cdbfc6"}, {"_id": "5e5218466a4cea0021cdbfc7"}, {"_id": "5e5218466a4cea0021cdbfc8"}, {"_id": "5e5218466a4cea0021cdbfc9"}, {"_id": "5e5218466a4cea0021cdbfca"}, {"_id": "5e5218466a4cea0021cdbfcb"}, {"_id": "5e5218466a4cea0021cdbfcc"}, {"_id": "5e5218466a4cea0021cdbfcd"}, {"_id": "5e5218466a4cea0021cdbfce"}, {"_id": "5e5218466a4cea0021cdbfcf"}, {"_id": "5e5218466a4cea0021cdbfd0"}, {"_id": "5e5218466a4cea0021cdbfd1"}, {"_id": "5e5218466a4cea0021cdbfd2"}, {"_id": "5e5218466a4cea0021cdbfd3"}, {"_id": "5e5218476a4cea0021cdbfd4"}, {"_id": "5e5218476a4cea0021cdbfd5"}, {"_id": "5e5218476a4cea0021cdbfd6"}, {"_id": "5e5218476a4cea0021cdbfd7"}, {"_id": "5e5218476a4cea0021cdbfd8"}, {"_id": "5e5218476a4cea0021cdbfd9"}, {"_id": "5e5218476a4cea0021cdbfda"}, {"_id": "5e5218476a4cea0021cdbfdb"}, {"_id": "5e5218476a4cea0021cdbfdc"}, {"_id": "5e5218476a4cea0021cdbfdd"}, {"_id": "5e5218476a4cea0021cdbfde"}, {"_id": "5e5218476a4cea0021cdbfdf"}, {"_id": "5e5218476a4cea0021cdbfe0"}, {"_id": "5e5218476a4cea0021cdbfe1"}, {"_id": "5e5218476a4cea0021cdbfe2"}, {"_id": "5e5218476a4cea0021cdbfe3"}, {"_id": "5e5218486a4cea0021cdbfe4"}, {"_id": "5e5218486a4cea0021cdbfe5"}, {"_id": "5e5218486a4cea0021cdbfe6"}, {"_id": "5e5218486a4cea0021cdbfe7"}, {"_id": "5e5218486a4cea0021cdbfe8"}, {"_id": "5e5218486a4cea0021cdbfe9"}, {"_id": "5e5218486a4cea0021cdbfea"}, {"_id": "5e5218486a4cea0021cdbfeb"}, {"_id": "5e5218486a4cea0021cdbfec"}, {"_id": "5e5218486a4cea0021cdbfed"}, {"_id": "5e5218486a4cea0021cdbfee"}, {"_id": "5e5218486a4cea0021cdbfef"}, {"_id": "5e5218486a4cea0021cdbff0"}, {"_id": "5e5218486a4cea0021cdbff1"}, {"_id": "5e5218486a4cea0021cdbff2"}, {"_id": "5e5218496a4cea0021cdbff3"}, {"_id": "5e5218496a4cea0021cdbff4"}, {"_id": "5e5218496a4cea0021cdbff5"}, {"_id": "5e5218496a4cea0021cdbff6"}, {"_id": "5e5218496a4cea0021cdbff7"}, {"_id": "5e5218496a4cea0021cdbff8"}, {"_id": "5e5218496a4cea0021cdbff9"}, {"_id": "5e5218496a4cea0021cdbffa"}, {"_id": "5e5218496a4cea0021cdbffb"}, {"_id": "5e5218496a4cea0021cdbffc"}, {"_id": "5e5218496a4cea0021cdbffd"}, {"_id": "5e5218496a4cea0021cdbffe"}, {"_id": "5e5218496a4cea0021cdbfff"}, {"_id": "5e5218496a4cea0021cdc000"}, {"_id": "5e5218496a4cea0021cdc001"}, {"_id": "5e52184a6a4cea0021cdc002"}, {"_id": "5e52184a6a4cea0021cdc003"}, {"_id": "5e52184a6a4cea0021cdc004"}, {"_id": "5e52184a6a4cea0021cdc005"}, {"_id": "5e52184a6a4cea0021cdc006"}, {"_id": "5e52184a6a4cea0021cdc007"}, {"_id": "5e52184a6a4cea0021cdc008"}, {"_id": "5e52184a6a4cea0021cdc009"}, {"_id": "5e52184a6a4cea0021cdc00a"}, {"_id": "5e52184a6a4cea0021cdc00b"}, {"_id": "5e52184a6a4cea0021cdc00c"}, {"_id": "5e52184a6a4cea0021cdc00d"}, {"_id": "5e52184a6a4cea0021cdc00e"}, {"_id": "5e52184a6a4cea0021cdc00f"}, {"_id": "5e52184b6a4cea0021cdc010"}, {"_id": "5e52184b6a4cea0021cdc011"}, {"_id": "5e52184b6a4cea0021cdc012"}, {"_id": "5e52184b6a4cea0021cdc013"}, {"_id": "5e52184b6a4cea0021cdc014"}, {"_id": "5e52184b6a4cea0021cdc015"}, {"_id": "5e52184b6a4cea0021cdc016"}, {"_id": "5e52184b6a4cea0021cdc017"}, {"_id": "5e52184b6a4cea0021cdc018"}, {"_id": "5e52184b6a4cea0021cdc019"}, {"_id": "5e52184b6a4cea0021cdc01a"}, {"_id": "5e52184b6a4cea0021cdc01b"}, {"_id": "5e52184b6a4cea0021cdc01c"}, {"_id": "5e52184b6a4cea0021cdc01d"}, {"_id": "5e52184b6a4cea0021cdc01e"}, {"_id": "5e52184c6a4cea0021cdc01f"}, {"_id": "5e52184c6a4cea0021cdc020"}, {"_id": "5e52184c6a4cea0021cdc021"}, {"_id": "5e52184c6a4cea0021cdc022"}, {"_id": "5e52184c6a4cea0021cdc023"}, {"_id": "5e52184c6a4cea0021cdc024"}, {"_id": "5e52184c6a4cea0021cdc025"}, {"_id": "5e52184c6a4cea0021cdc026"}, {"_id": "5e52184c6a4cea0021cdc027"}, {"_id": "5e52184c6a4cea0021cdc028"}, {"_id": "5e52184c6a4cea0021cdc029"}, {"_id": "5e52184c6a4cea0021cdc02a"}, {"_id": "5e52184c6a4cea0021cdc02b"}, {"_id": "5e52184c6a4cea0021cdc02c"}, {"_id": "5e52184c6a4cea0021cdc02d"}, {"_id": "5e52184c6a4cea0021cdc02e"}, {"_id": "5e52184d6a4cea0021cdc02f"}, {"_id": "5e52184d6a4cea0021cdc030"}, {"_id": "5e52184d6a4cea0021cdc031"}, {"_id": "5e52184d6a4cea0021cdc032"}, {"_id": "5e52184d6a4cea0021cdc033"}, {"_id": "5e52184d6a4cea0021cdc034"}, {"_id": "5e52184d6a4cea0021cdc035"}, {"_id": "5e52184d6a4cea0021cdc036"}, {"_id": "5e52184d6a4cea0021cdc037"}, {"_id": "5e52184d6a4cea0021cdc038"}, {"_id": "5e52184d6a4cea0021cdc039"}, {"_id": "5e52184d6a4cea0021cdc03a"}, {"_id": "5e52184d6a4cea0021cdc03b"}, {"_id": "5e52184d6a4cea0021cdc03c"}, {"_id": "5e52184d6a4cea0021cdc03d"}, {"_id": "5e52184d6a4cea0021cdc03e"}, {"_id": "5e52184e6a4cea0021cdc03f"}, {"_id": "5e52184e6a4cea0021cdc040"}, {"_id": "5e52184e6a4cea0021cdc041"}, {"_id": "5e52184e6a4cea0021cdc042"}, {"_id": "5e52184e6a4cea0021cdc043"}, {"_id": "5e52184e6a4cea0021cdc044"}, {"_id": "5e52184e6a4cea0021cdc045"}, {"_id": "5e52184e6a4cea0021cdc046"}, {"_id": "5e52184e6a4cea0021cdc047"}, {"_id": "5e52184e6a4cea0021cdc048"}, {"_id": "5e52184e6a4cea0021cdc049"}, {"_id": "5e52184e6a4cea0021cdc04a"}, {"_id": "5e52184e6a4cea0021cdc04b"}, {"_id": "5e52184e6a4cea0021cdc04c"}, {"_id": "5e52184e6a4cea0021cdc04d"}, {"_id": "5e52184e6a4cea0021cdc04e"}, {"_id": "5e52184f6a4cea0021cdc04f"}, {"_id": "5e52184f6a4cea0021cdc050"}, {"_id": "5e52184f6a4cea0021cdc051"}, {"_id": "5e52184f6a4cea0021cdc052"}, {"_id": "5e52184f6a4cea0021cdc053"}, {"_id": "5e52184f6a4cea0021cdc054"}, {"_id": "5e52184f6a4cea0021cdc055"}, {"_id": "5e52184f6a4cea0021cdc056"}, {"_id": "5e52184f6a4cea0021cdc057"}, {"_id": "5e52184f6a4cea0021cdc058"}, {"_id": "5e52184f6a4cea0021cdc059"}, {"_id": "5e52184f6a4cea0021cdc05a"}, {"_id": "5e52184f6a4cea0021cdc05b"}, {"_id": "5e52184f6a4cea0021cdc05c"}, {"_id": "5e52184f6a4cea0021cdc05d"}, {"_id": "5e52184f6a4cea0021cdc05e"}, {"_id": "5e5218506a4cea0021cdc05f"}, {"_id": "5e5218506a4cea0021cdc060"}, {"_id": "5e5218506a4cea0021cdc061"}, {"_id": "5e5218506a4cea0021cdc062"}, {"_id": "5e5218506a4cea0021cdc063"}, {"_id": "5e5218506a4cea0021cdc064"}, {"_id": "5e5218506a4cea0021cdc065"}, {"_id": "5e5218506a4cea0021cdc066"}, {"_id": "5e5218506a4cea0021cdc067"}, {"_id": "5e5218506a4cea0021cdc068"}, {"_id": "5e5218506a4cea0021cdc069"}, {"_id": "5e5218506a4cea0021cdc06a"}, {"_id": "5e5218506a4cea0021cdc06b"}, {"_id": "5e5218506a4cea0021cdc06c"}, {"_id": "5e5218506a4cea0021cdc06d"}, {"_id": "5e5218506a4cea0021cdc06e"}, {"_id": "5e5218516a4cea0021cdc06f"}, {"_id": "5e5218516a4cea0021cdc070"}, {"_id": "5e5218516a4cea0021cdc071"}, {"_id": "5e5218516a4cea0021cdc072"}, {"_id": "5e5218516a4cea0021cdc073"}, {"_id": "5e5218516a4cea0021cdc074"}, {"_id": "5e5218516a4cea0021cdc075"}, {"_id": "5e5218516a4cea0021cdc076"}, {"_id": "5e5218516a4cea0021cdc077"}, {"_id": "5e5218516a4cea0021cdc078"}, {"_id": "5e5218516a4cea0021cdc079"}, {"_id": "5e5218516a4cea0021cdc07a"}, {"_id": "5e5218516a4cea0021cdc07b"}, {"_id": "5e5218516a4cea0021cdc07c"}, {"_id": "5e5218516a4cea0021cdc07d"}, {"_id": "5e5218516a4cea0021cdc07e"}, {"_id": "5e5218516a4cea0021cdc07f"}, {"_id": "5e5218526a4cea0021cdc080"}, {"_id": "5e5218526a4cea0021cdc081"}, {"_id": "5e5218526a4cea0021cdc082"}, {"_id": "5e5218526a4cea0021cdc083"}, {"_id": "5e5218526a4cea0021cdc084"}, {"_id": "5e5218526a4cea0021cdc085"}, {"_id": "5e5218526a4cea0021cdc086"}, {"_id": "5e5218526a4cea0021cdc087"}, {"_id": "5e5218526a4cea0021cdc088"}, {"_id": "5e5218526a4cea0021cdc089"}, {"_id": "5e5218526a4cea0021cdc08a"}, {"_id": "5e5218526a4cea0021cdc08b"}, {"_id": "5e5218526a4cea0021cdc08c"}, {"_id": "5e5218526a4cea0021cdc08d"}, {"_id": "5e5218526a4cea0021cdc08e"}, {"_id": "5e5218526a4cea0021cdc08f"}, {"_id": "5e5218526a4cea0021cdc090"}, {"_id": "5e5218536a4cea0021cdc091"}, {"_id": "5e5218536a4cea0021cdc092"}, {"_id": "5e5218536a4cea0021cdc093"}, {"_id": "5e5218536a4cea0021cdc094"}, {"_id": "5e5218536a4cea0021cdc095"}, {"_id": "5e5218536a4cea0021cdc096"}, {"_id": "5e5218536a4cea0021cdc097"}, {"_id": "5e5218536a4cea0021cdc098"}, {"_id": "5e5218536a4cea0021cdc099"}, {"_id": "5e5218536a4cea0021cdc09a"}, {"_id": "5e5218536a4cea0021cdc09b"}, {"_id": "5e5218536a4cea0021cdc09c"}, {"_id": "5e5218536a4cea0021cdc09d"}, {"_id": "5e5218536a4cea0021cdc09e"}, {"_id": "5e5218546a4cea0021cdc09f"}, {"_id": "5e5218546a4cea0021cdc0a0"}, {"_id": "5e5218546a4cea0021cdc0a1"}, {"_id": "5e5218546a4cea0021cdc0a2"}, {"_id": "5e5218546a4cea0021cdc0a3"}, {"_id": "5e5218546a4cea0021cdc0a4"}, {"_id": "5e5218546a4cea0021cdc0a5"}, {"_id": "5e5218546a4cea0021cdc0a6"}, {"_id": "5e5218546a4cea0021cdc0a7"}, {"_id": "5e5218546a4cea0021cdc0a8"}, {"_id": "5e5218546a4cea0021cdc0a9"}, {"_id": "5e5218546a4cea0021cdc0aa"}, {"_id": "5e5218546a4cea0021cdc0ab"}, {"_id": "5e5218546a4cea0021cdc0ac"}, {"_id": "5e5218546a4cea0021cdc0ad"}, {"_id": "5e5218546a4cea0021cdc0ae"}, {"_id": "5e5218546a4cea0021cdc0af"}, {"_id": "5e5218556a4cea0021cdc0b0"}, {"_id": "5e5218556a4cea0021cdc0b1"}, {"_id": "5e5218556a4cea0021cdc0b2"}, {"_id": "5e5218556a4cea0021cdc0b3"}, {"_id": "5e5218556a4cea0021cdc0b4"}, {"_id": "5e5218556a4cea0021cdc0b5"}, {"_id": "5e5218556a4cea0021cdc0b6"}, {"_id": "5e5218556a4cea0021cdc0b7"}, {"_id": "5e5218556a4cea0021cdc0b8"}, {"_id": "5e5218556a4cea0021cdc0b9"}, {"_id": "5e5218556a4cea0021cdc0ba"}, {"_id": "5e5218556a4cea0021cdc0bb"}, {"_id": "5e5218556a4cea0021cdc0bc"}, {"_id": "5e5218556a4cea0021cdc0bd"}, {"_id": "5e5218556a4cea0021cdc0be"}, {"_id": "5e5218556a4cea0021cdc0bf"}, {"_id": "5e5218556a4cea0021cdc0c0"}, {"_id": "5e5218566a4cea0021cdc0c1"}, {"_id": "5e5218566a4cea0021cdc0c2"}, {"_id": "5e5218566a4cea0021cdc0c3"}, {"_id": "5e5218566a4cea0021cdc0c4"}, {"_id": "5e5218566a4cea0021cdc0c5"}, {"_id": "5e5218566a4cea0021cdc0c6"}, {"_id": "5e5218566a4cea0021cdc0c7"}, {"_id": "5e5218566a4cea0021cdc0c8"}, {"_id": "5e5218566a4cea0021cdc0c9"}, {"_id": "5e5218566a4cea0021cdc0ca"}, {"_id": "5e5218566a4cea0021cdc0cb"}, {"_id": "5e5218566a4cea0021cdc0cc"}, {"_id": "5e5218566a4cea0021cdc0cd"}, {"_id": "5e5218566a4cea0021cdc0ce"}, {"_id": "5e5218566a4cea0021cdc0cf"}, {"_id": "5e5218566a4cea0021cdc0d0"}, {"_id": "5e5218566a4cea0021cdc0d1"}, {"_id": "5e5218576a4cea0021cdc0d2"}, {"_id": "5e5218576a4cea0021cdc0d3"}, {"_id": "5e5218576a4cea0021cdc0d4"}, {"_id": "5e5218576a4cea0021cdc0d5"}, {"_id": "5e5218576a4cea0021cdc0d6"}, {"_id": "5e5218576a4cea0021cdc0d7"}, {"_id": "5e5218576a4cea0021cdc0d8"}, {"_id": "5e5218576a4cea0021cdc0d9"}, {"_id": "5e5218576a4cea0021cdc0da"}, {"_id": "5e5218576a4cea0021cdc0db"}, {"_id": "5e5218576a4cea0021cdc0dc"}, {"_id": "5e5218576a4cea0021cdc0dd"}, {"_id": "5e5218586a4cea0021cdc0de"}, {"_id": "5e5218586a4cea0021cdc0df"}, {"_id": "5e5218586a4cea0021cdc0e0"}, {"_id": "5e5218586a4cea0021cdc0e1"}, {"_id": "5e5218586a4cea0021cdc0e2"}, {"_id": "5e5218586a4cea0021cdc0e3"}, {"_id": "5e5218586a4cea0021cdc0e4"}, {"_id": "5e5218586a4cea0021cdc0e5"}, {"_id": "5e5218586a4cea0021cdc0e6"}, {"_id": "5e5218586a4cea0021cdc0e7"}, {"_id": "5e5218586a4cea0021cdc0e8"}, {"_id": "5e5218586a4cea0021cdc0e9"}, {"_id": "5e5218586a4cea0021cdc0ea"}, {"_id": "5e5218586a4cea0021cdc0eb"}, {"_id": "5e5218586a4cea0021cdc0ec"}, {"_id": "5e5218596a4cea0021cdc0ed"}, {"_id": "5e5218596a4cea0021cdc0ee"}, {"_id": "5e5218596a4cea0021cdc0ef"}, {"_id": "5e5218596a4cea0021cdc0f0"}, {"_id": "5e5218596a4cea0021cdc0f1"}, {"_id": "5e5218596a4cea0021cdc0f2"}, {"_id": "5e5218596a4cea0021cdc0f3"}, {"_id": "5e5218596a4cea0021cdc0f4"}, {"_id": "5e5218596a4cea0021cdc0f5"}, {"_id": "5e5218596a4cea0021cdc0f6"}, {"_id": "5e5218596a4cea0021cdc0f7"}, {"_id": "5e5218596a4cea0021cdc0f8"}, {"_id": "5e5218596a4cea0021cdc0f9"}, {"_id": "5e5218596a4cea0021cdc0fa"}, {"_id": "5e5218596a4cea0021cdc0fb"}, {"_id": "5e5218596a4cea0021cdc0fc"}, {"_id": "5e5218596a4cea0021cdc0fd"}, {"_id": "5e52185a6a4cea0021cdc0fe"}, {"_id": "5e52185a6a4cea0021cdc0ff"}, {"_id": "5e52185a6a4cea0021cdc100"}, {"_id": "5e52185a6a4cea0021cdc101"}, {"_id": "5e52185a6a4cea0021cdc102"}, {"_id": "5e52185a6a4cea0021cdc103"}, {"_id": "5e52185a6a4cea0021cdc104"}, {"_id": "5e52185a6a4cea0021cdc105"}, {"_id": "5e52185a6a4cea0021cdc106"}, {"_id": "5e52185a6a4cea0021cdc107"}, {"_id": "5e52185a6a4cea0021cdc108"}, {"_id": "5e52185a6a4cea0021cdc109"}, {"_id": "5e52185a6a4cea0021cdc10a"}, {"_id": "5e52185a6a4cea0021cdc10b"}, {"_id": "5e52185a6a4cea0021cdc10c"}, {"_id": "5e52185a6a4cea0021cdc10d"}, {"_id": "5e52185a6a4cea0021cdc10e"}, {"_id": "5e52185b6a4cea0021cdc10f"}, {"_id": "5e52185b6a4cea0021cdc110"}, {"_id": "5e52185b6a4cea0021cdc111"}, {"_id": "5e52185b6a4cea0021cdc112"}, {"_id": "5e52185b6a4cea0021cdc113"}, {"_id": "5e52185b6a4cea0021cdc114"}, {"_id": "5e52185b6a4cea0021cdc115"}, {"_id": "5e52185b6a4cea0021cdc116"}, {"_id": "5e52185b6a4cea0021cdc117"}, {"_id": "5e52185b6a4cea0021cdc118"}, {"_id": "5e52185b6a4cea0021cdc119"}, {"_id": "5e52185b6a4cea0021cdc11a"}, {"_id": "5e52185b6a4cea0021cdc11b"}, {"_id": "5e52185b6a4cea0021cdc11c"}, {"_id": "5e52185b6a4cea0021cdc11d"}, {"_id": "5e52185c6a4cea0021cdc11e"}, {"_id": "5e52185c6a4cea0021cdc11f"}, {"_id": "5e52185c6a4cea0021cdc120"}, {"_id": "5e52185c6a4cea0021cdc121"}, {"_id": "5e52185c6a4cea0021cdc122"}, {"_id": "5e52185c6a4cea0021cdc123"}, {"_id": "5e52185c6a4cea0021cdc124"}, {"_id": "5e52185c6a4cea0021cdc125"}, {"_id": "5e52185c6a4cea0021cdc126"}, {"_id": "5e52185c6a4cea0021cdc127"}, {"_id": "5e52185c6a4cea0021cdc128"}, {"_id": "5e52185c6a4cea0021cdc129"}, {"_id": "5e52185c6a4cea0021cdc12a"}, {"_id": "5e52185c6a4cea0021cdc12b"}, {"_id": "5e52185c6a4cea0021cdc12c"}, {"_id": "5e52185c6a4cea0021cdc12d"}, {"_id": "5e52185d6a4cea0021cdc12e"}, {"_id": "5e52185d6a4cea0021cdc12f"}, {"_id": "5e52185d6a4cea0021cdc130"}, {"_id": "5e52185d6a4cea0021cdc131"}, {"_id": "5e52185d6a4cea0021cdc132"}, {"_id": "5e52185d6a4cea0021cdc133"}, {"_id": "5e52185d6a4cea0021cdc134"}, {"_id": "5e52185d6a4cea0021cdc135"}, {"_id": "5e52185d6a4cea0021cdc136"}, {"_id": "5e52185d6a4cea0021cdc137"}, {"_id": "5e52185d6a4cea0021cdc138"}, {"_id": "5e52185d6a4cea0021cdc139"}, {"_id": "5e52185d6a4cea0021cdc13a"}, {"_id": "5e52185d6a4cea0021cdc13b"}, {"_id": "5e52185d6a4cea0021cdc13c"}, {"_id": "5e52185d6a4cea0021cdc13d"}, {"_id": "5e52185d6a4cea0021cdc13e"}, {"_id": "5e52185e6a4cea0021cdc13f"}, {"_id": "5e52185e6a4cea0021cdc140"}, {"_id": "5e52185e6a4cea0021cdc141"}, {"_id": "5e52185e6a4cea0021cdc142"}, {"_id": "5e52185e6a4cea0021cdc143"}, {"_id": "5e52185e6a4cea0021cdc144"}, {"_id": "5e52185e6a4cea0021cdc145"}, {"_id": "5e52185e6a4cea0021cdc146"}, {"_id": "5e52185e6a4cea0021cdc147"}, {"_id": "5e52185e6a4cea0021cdc148"}, {"_id": "5e52185e6a4cea0021cdc149"}, {"_id": "5e52185e6a4cea0021cdc14a"}, {"_id": "5e52185e6a4cea0021cdc14b"}, {"_id": "5e52185e6a4cea0021cdc14c"}, {"_id": "5e52185e6a4cea0021cdc14d"}, {"_id": "5e52185e6a4cea0021cdc14e"}, {"_id": "5e52185f6a4cea0021cdc14f"}, {"_id": "5e52185f6a4cea0021cdc150"}, {"_id": "5e52185f6a4cea0021cdc151"}, {"_id": "5e52185f6a4cea0021cdc152"}, {"_id": "5e52185f6a4cea0021cdc153"}, {"_id": "5e52185f6a4cea0021cdc154"}, {"_id": "5e52185f6a4cea0021cdc155"}, {"_id": "5e52185f6a4cea0021cdc156"}, {"_id": "5e52185f6a4cea0021cdc157"}, {"_id": "5e52185f6a4cea0021cdc158"}, {"_id": "5e52185f6a4cea0021cdc159"}, {"_id": "5e52185f6a4cea0021cdc15a"}, {"_id": "5e52185f6a4cea0021cdc15b"}, {"_id": "5e52185f6a4cea0021cdc15c"}, {"_id": "5e52185f6a4cea0021cdc15d"}, {"_id": "5e5218606a4cea0021cdc15e"}, {"_id": "5e5218606a4cea0021cdc15f"}, {"_id": "5e5218606a4cea0021cdc160"}, {"_id": "5e5218606a4cea0021cdc161"}, {"_id": "5e5218606a4cea0021cdc162"}, {"_id": "5e5218606a4cea0021cdc163"}, {"_id": "5e5218606a4cea0021cdc164"}, {"_id": "5e5218606a4cea0021cdc165"}, {"_id": "5e5218606a4cea0021cdc166"}, {"_id": "5e5218606a4cea0021cdc167"}, {"_id": "5e5218606a4cea0021cdc168"}, {"_id": "5e5218606a4cea0021cdc169"}, {"_id": "5e5218606a4cea0021cdc16a"}, {"_id": "5e5218606a4cea0021cdc16b"}, {"_id": "5e5218606a4cea0021cdc16c"}, {"_id": "5e5218606a4cea0021cdc16d"}, {"_id": "5e5218606a4cea0021cdc16e"}, {"_id": "5e5218616a4cea0021cdc16f"}, {"_id": "5e5218616a4cea0021cdc170"}, {"_id": "5e5218616a4cea0021cdc171"}, {"_id": "5e5218616a4cea0021cdc172"}, {"_id": "5e5218616a4cea0021cdc173"}, {"_id": "5e5218616a4cea0021cdc174"}, {"_id": "5e5218616a4cea0021cdc175"}, {"_id": "5e5218616a4cea0021cdc176"}, {"_id": "5e5218616a4cea0021cdc177"}, {"_id": "5e5218616a4cea0021cdc178"}, {"_id": "5e5218616a4cea0021cdc179"}, {"_id": "5e5218616a4cea0021cdc17a"}, {"_id": "5e5218616a4cea0021cdc17b"}, {"_id": "5e5218616a4cea0021cdc17c"}, {"_id": "5e5218616a4cea0021cdc17d"}, {"_id": "5e5218616a4cea0021cdc17e"}, {"_id": "5e5218616a4cea0021cdc17f"}, {"_id": "5e5218626a4cea0021cdc180"}, {"_id": "5e5218626a4cea0021cdc181"}, {"_id": "5e5218626a4cea0021cdc182"}, {"_id": "5e5218626a4cea0021cdc183"}, {"_id": "5e5218626a4cea0021cdc184"}, {"_id": "5e5218626a4cea0021cdc185"}, {"_id": "5e5218626a4cea0021cdc186"}, {"_id": "5e5218626a4cea0021cdc187"}, {"_id": "5e5218626a4cea0021cdc188"}, {"_id": "5e5218626a4cea0021cdc189"}, {"_id": "5e5218626a4cea0021cdc18a"}, {"_id": "5e5218626a4cea0021cdc18b"}, {"_id": "5e5218626a4cea0021cdc18c"}, {"_id": "5e5218626a4cea0021cdc18d"}];
    const AuthedUsers = [{
        "name": "3d",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOGY4Njc0MzhkYTAwMjBmNTZjOWQiLCJpYXQiOjE1ODMzNzAxNTB9.7qokxb8xZxsW2k3SF8H9Cxj5XGZBP2JUSpDKpxxS6kI"
    }, {
        "name": "a",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTA3OGNiMDc2YjAwMWYzYWJhMzkiLCJpYXQiOjE1ODMzNzAxNTF9.F7vq0ht0S2NijbjrjgugygUme0Tm-_yx2GrhKlNblGI"
    }, {
        "name": "aa",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTA3ODc0MzhkYTAwMjBmNTZjYTIiLCJpYXQiOjE1ODMzNzAxNTF9.RIUDJ0WSq-Crvq5ly9ltO47rzFh--Z-pNU20zbHoH3g"
    }, {
        "name": "aaberg",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTA3OGNiMDc2YjAwMWYzYWJhM2EiLCJpYXQiOjE1ODMzNzAxNTF9.LySYXKmDavsBri_Q0cd0Nn8gNpyap_LzisR2ZVzeJ3I"
    }, {
        "name": "a1",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTA3OGU4OTVjODAwMThjNjBkMzQiLCJpYXQiOjE1ODMzNzAxNTF9.inZwQwSapSgLU9ZtP4q8YtKptbtICJV2XMoSL08JNNM"
    }, {
        "name": "aachen",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkMzc0MzhkYTAwMjBmNTZjYWEiLCJpYXQiOjE1ODMzNzAxNTF9.d0_uptTTnQk0QT0EVJ0E-ARKw4-wtSht9IS1u2sORVU"
    }, {
        "name": "aalborg",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2U4OTVjODAwMThjNjBkM2QiLCJpYXQiOjE1ODMzNzAxNTF9.dVkyB5Ap-vHgJkKkcegKjOiLH760aujGnaoAREiuu9c"
    }, {
        "name": "aalii",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2NiMDc2YjAwMWYzYWJhNDQiLCJpYXQiOjE1ODMzNzAxNTF9.giTSjhG-_JcvBZJSeQfD5ilMOBJyDWXW_23gptGozKk"
    }, {
        "name": "aalesund",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2U4OTVjODAwMThjNjBkM2UiLCJpYXQiOjE1ODMzNzAxNTF9.QX3gSJ2WKe_lrzMyI7LgMfKiR3SCSDhWCU-U1kON1e4"
    }, {
        "name": "aalto",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkMzc0MzhkYTAwMjBmNTZjYWMiLCJpYXQiOjE1ODMzNzAxNTF9.H9w4BMtY4ZlEMauNDaicXH3CyN4t3FObhW2mnzX9w1g"
    }, {
        "name": "aalst",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2NiMDc2YjAwMWYzYWJhNDUiLCJpYXQiOjE1ODMzNzAxNTF9.4xkbhSUoUygb9sWZ-tmTLLkHXhw265I-YnRGNaR7K7U"
    }, {
        "name": "aam",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2U4OTVjODAwMThjNjBkNDEiLCJpYXQiOjE1ODMzNzAxNTF9.VA8Jn4pR_aUe-cTuiMgidKkSosuOfPhsOQls-45RAn0"
    }, {
        "name": "aarau",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkMzc0MzhkYTAwMjBmNTZjYWUiLCJpYXQiOjE1ODMzNzAxNTF9.inH-cKT7zbQBmp_rU89nE-j3ADapfH1lhDveGMqi8jE"
    }, {
        "name": "aardvark",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2NiMDc2YjAwMWYzYWJhNDgiLCJpYXQiOjE1ODMzNzAxNTJ9.lhyJUW2agso3SAmlXiDnoa3fuKdMDARmCcHV5B4WuTE"
    }, {
        "name": "aara",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2U4OTVjODAwMThjNjBkNDIiLCJpYXQiOjE1ODMzNzAxNTF9.wfnxhcmkf7qBIRXAK4eBPnLNVrFinJNwRL0ehMj6wbE"
    }, {
        "name": "aardwolf",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2NiMDc2YjAwMWYzYWJhNDkiLCJpYXQiOjE1ODMzNzAxNTJ9.lLNPMMzoT-kp5TfnPgZlzprDaGzuzaDR4lyL_8ZTJO8"
    }, {
        "name": "aaren",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkMzc0MzhkYTAwMjBmNTZjYjAiLCJpYXQiOjE1ODMzNzAxNTJ9.GSajCyCgk7HfwbiyUeQrq1O_5bGws7dc-wklIJQMCMc"
    }, {
        "name": "aargau",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkM2U4OTVjODAwMThjNjBkNDUiLCJpYXQiOjE1ODMzNzAxNTJ9.Ar5hqZgHYCAP8GMBQjyL_atyow-jcgYJMJpd_5tCvNc"
    }, {
        "name": "aarika",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNGMiLCJpYXQiOjE1ODMzNzAxNTJ9.9HTZFr0WmRol0sSgT_JPzDC3pdjjh547yI4gH1tNHfM"
    }, {
        "name": "aarhus",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNDYiLCJpYXQiOjE1ODMzNzAxNTJ9.VJC-LcUWspCpqpvW-pCozwZQoL2Sq21d77uSuawc_Ss"
    }, {
        "name": "aaron",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYjIiLCJpYXQiOjE1ODMzNzAxNTJ9.zCW7reiMrBQcWQ8HOzHvHo6_iSI6BSRg1PZVMN4rbYQ"
    }, {
        "name": "aaronaaronic",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNDkiLCJpYXQiOjE1ODMzNzAxNTJ9.HLkLAWKnPVhBVhP5C3PM4kt2keXePaV4S49g6VAI_W0"
    }, {
        "name": "ab",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYjQiLCJpYXQiOjE1ODMzNzAxNTJ9.5PvPXHVDslvAUcqweajY2YEWBw5QPKjTUxCzkGG5dGQ"
    }, {
        "name": "aba",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNGYiLCJpYXQiOjE1ODMzNzAxNTJ9.rh4R-nsGwt7mfhqim0idOyRFXWgXMRj-1ydrfevtE2M"
    }, {
        "name": "aaronson",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNGEiLCJpYXQiOjE1ODMzNzAxNTJ9.aXCqh-NrzTg0bw_54Vmw9O6Eva0M2UpUrhJ7VC9LgMs"
    }, {
        "name": "abacist",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYjYiLCJpYXQiOjE1ODMzNzAxNTN9.PVRCPqWlewh2S0WElnB_Lzpp8zUa9rd9AB4bgoRGFAE"
    }, {
        "name": "aback",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNGQiLCJpYXQiOjE1ODMzNzAxNTN9.AvT8PiJVEtoTy6UVX2l6vci7I4kG41tKSpmBmAnTLxo"
    }, {
        "name": "abaca",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNTAiLCJpYXQiOjE1ODMzNzAxNTN9.zwGy4hSTXGYmTtsNiI0lQ2nW5RYXBF2r26y1IJ-9TM8"
    }, {
        "name": "abacus",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNGUiLCJpYXQiOjE1ODMzNzAxNTN9.PWEzp6Y6bKDEvucpjJni6goBuGwLxb89Y8u29ItoJqs"
    }, {
        "name": "abad",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYjgiLCJpYXQiOjE1ODMzNzAxNTN9.dL14-kjCMPoaLln4a0ENQLAXwo0SlUDHFVVy57Xf7NY"
    }, {
        "name": "abaddon",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNTEiLCJpYXQiOjE1ODMzNzAxNTN9.dMbL3LXuU-7VfRrMDnsq-O59UZP30Qh-DnrGGPZ8iiw"
    }, {
        "name": "abadan",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNTQiLCJpYXQiOjE1ODMzNzAxNTN9.dQpZUC0MetqADpA8gJGe-YFC1H18H1R0z17ri9gbu4s"
    }, {
        "name": "abagael",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYmEiLCJpYXQiOjE1ODMzNzAxNTN9.jQ_PpOzXol34eTuJAEuW_HpqLh4KrRqGiYg7QgCkFas"
    }, {
        "name": "abaft",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNTIiLCJpYXQiOjE1ODMzNzAxNTN9.UavlE0lHUd387iC49l2BIeBfIPo5QAATwmnN87a9FFI"
    }, {
        "name": "abagail",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNTYiLCJpYXQiOjE1ODMzNzAxNTN9.ZLZXO6NeX0gHXyJUxD3wE7qzJeNmqD6c6ng9aYSRMGs"
    }, {
        "name": "abamp",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYmMiLCJpYXQiOjE1ODMzNzAxNTN9.N9g-rKrUs7JDPWlAiB8X_T8u86Fzug_BUy5fDj6f7aQ"
    }, {
        "name": "abampere",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNTUiLCJpYXQiOjE1ODMzNzAxNTR9.G8yYFODnkVmQzVTEAKZ2VMGaMun-sKnnwRzAkEpHYKk"
    }, {
        "name": "abalone",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNTciLCJpYXQiOjE1ODMzNzAxNTN9.5gy00-qf2DZvbxV4T3s2Zjqd9IPa2aSPzchoInYZLUA"
    }, {
        "name": "abandon",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYmUiLCJpYXQiOjE1ODMzNzAxNTR9.vc3rTTbnvgwmWYJ1XI6bChOvMUbL85GUrAVau4uziQk"
    }, {
        "name": "abana",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNTYiLCJpYXQiOjE1ODMzNzAxNTR9.wWsyoySkUgNgFKGhnwMIH-9dNoTJ7B_W6zNkaCfvBDg"
    }, {
        "name": "abandoned",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNWEiLCJpYXQiOjE1ODMzNzAxNTR9.S-lxnpYXoDlMWkzKJ-dVAeEAif9yEiHdamtE4XwKWic"
    }, {
        "name": "abase",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGU4OTVjODAwMThjNjBkNTkiLCJpYXQiOjE1ODMzNzAxNTR9.sK6JydItvd6JbpM_gVQUQLum3TDHSkcvG-iOLjvxY4U"
    }, {
        "name": "abash",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNDc0MzhkYTAwMjBmNTZjYzAiLCJpYXQiOjE1ODMzNzAxNTR9.v9rFhooxBiVKOA3nHMoLLDaCYsInYXgYslc51g9Bizo"
    }, {
        "name": "abarca",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNGNiMDc2YjAwMWYzYWJhNWIiLCJpYXQiOjE1ODMzNzAxNTR9.yH7n6dCiRwq_QH8Neicah5PnAI5-s-03z9YMv3ddjec"
    }, {
        "name": "abate",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNWU4OTVjODAwMThjNjBkNWEiLCJpYXQiOjE1ODMzNzAxNTR9.M8bEFAYKBxwfVZQz4aZh3Ne5v2Me-zrNRL1Pehdi93Y"
    }, {
        "name": "abatement",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNTc0MzhkYTAwMjBmNTZjYzIiLCJpYXQiOjE1ODMzNzAxNTR9.DqimMc6f6ls-b4Tof2Bmri9pwArbZBvHPWrf2IrX-KQ"
    }, {
        "name": "abatis",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNWNiMDc2YjAwMWYzYWJhNWUiLCJpYXQiOjE1ODMzNzAxNTR9.sfhzSl02VngdkPmY1WiffxvzYc51___zsJa2gApxX_w"
    }, {
        "name": "abattoir",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBkNWNiMDc2YjAwMWYzYWJhNWYiLCJpYXQiOjE1ODMzNzAxNTR9.8PZ6YJzwBEfmNkvKJUAeeUjul7K7pCrwidl63fIp874"
    }, {
        "name": "abb",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZTUxOTBmN2NiMDc2YjAwMWYzYWJhNzEiLCJpYXQiOjE1ODMzNzAxNTR9.bhehspyeDIMsbSdZ6Pl7VxbysazGCMArGPdEvrWc7M0"
    }];


    let uniqueNumber = __VU * 100000000 + (__ITER % 4000);
    randomSeed(uniqueNumber);

    const LoginAuthProb = argsloginprob;
    const GetBookProb = argsgetbookprob;
    const EditBookProb = argseditbookprob;
    const SLEEP_DURATION = argssleepduration;

    function get_random_item(items) {
        const item = items[Math.floor(Math.random() * items.length)];
        return item
    }

    const execute_edit_book = function () {
        const auth_token = get_random_item(AuthedUsers)['token'];
        const book = get_random_item(books);
        book['pages'] = Math.floor(Math.random() * 200) + 200;
        let edit_book_params = {
            headers: {
                'Authorization': 'Bearer ' + auth_token,
                'debug_id': new Date().getTime()
            },
            tags: {
                name: "edit_book"
            }
        };
        let edit_book_response = http.put(
            BaseURL + '/books/' + book._id,
            JSON.stringify(book),
            edit_book_params,
        );
        edit_book_duration_trend.add(edit_book_response.timings.duration);
        edit_book_counter.add(1);
        check(edit_book_response, {
            'is_edit_book_200': r => r.status === 200
        });
    }
    const execute_get_book = function () {
        const auth_token = get_random_item(AuthedUsers)['token'];
        const book = get_random_item(books);
        let get_book_params = {
            headers: {
                'Authorization': 'Bearer ' + auth_token,
                'debug_id': new Date().getTime()
            },
            tags: {
                name: "get_book"
            }
        };
        let get_book_response = http.get(
            BaseURL + '/books/' + book._id,
            get_book_params
        );
        get_book_duration_trend.add(get_book_response.timings.duration);
        get_book_counter.add(1);
        check(get_book_response, {
            'is_get_book_200': r => r.status === 200
        });
    }
    const execute_random_login = function () {
        const name = get_random_item(AuthedUsers)['name']
        const email = name + "@gmail.com";
        const password = "123456789";

        let body = JSON.stringify({
            email: email,
            password: password,
        });

        let login_params = {
            headers: {
                'Content-Type': 'application/json',
                'debug_id': new Date().getTime(),
            },
            tags: {
                name: "login"
            }
        };

        let login_response = http.post(
            BaseURL + '/auth/login',
            body,
            login_params
        );
        login_duration_trend.add(login_response.timings.duration);
        login_counter.add(1);
        check(login_response, {
            'is_login_200': r => r.status === 200,
            'is api key present': r => r.json().hasOwnProperty('token'),
        });
        AuthedUsers.push({name, "token": login_response.json()["token"]})
    }

    const r = Math.random();
    const sTime = Math.random() * SLEEP_DURATION + 0.5 * SLEEP_DURATION;
    sleep(sTime);
    if (r < LoginAuthProb) {
        execute_random_login();
    } else if (r >= LoginAuthProb && r < LoginAuthProb + GetBookProb) {
        execute_get_book();
    } else {
        execute_edit_book();
    }

};

export function teardown(data) {
}
