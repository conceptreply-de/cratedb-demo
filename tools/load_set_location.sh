#!/bin/bash

K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_PERIOD='1s' k6 run tools/load_set_location_k6.js