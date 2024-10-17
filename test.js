// Replace with your actual localhost URL and endpoint
const url = 'http://localhost:8080/test_db'; 

fetch(url, {
    method: 'GET', // Since it's a GET request, this is optional
    headers: {
        'Content-Type': 'application/json', // Expecting a JSON response
    }
})
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json(); // Parsing the JSON response
    })
    .then(data => {
        console.log('Success:', data); // Handle the JSON response
    })
    .catch((error) => {
        console.error('Error testing API:', error);
    });
