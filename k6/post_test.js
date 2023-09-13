import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10, // Virtual Users (concurrent connections)
    duration: '30s', // Test duration
};

export default function () {
    // Define your API endpoint
    let apiUrl = 'http://localhost:3000/todos'; // Update with your API URL

    // Define the request body (JSON format)
    let requestBody = JSON.stringify({
        "title": "Test 13",
        "description": "Todooooo!",
        "isCompleted": false
    });

    // Define request headers (if needed)
    let headers = {
        'Content-Type': 'application/json',
        // Add any other headers as needed
    };

    // Send an HTTP POST request with the request body and headers
    let response = http.post(apiUrl, requestBody, { headers: headers });

    // Check the response to ensure it's successful (e.g., status code 200)
    check(response, {
        'is status 200': (r) => r.status === 200,
    });

    // Add a sleep to simulate user think time (optional)
    sleep(1); // Sleep for 1 second
}
