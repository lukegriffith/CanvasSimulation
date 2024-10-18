const canvas = document.getElementById('simulationCanvas');
const ctx = canvas.getContext('2d');

// Connect to the WebSocket server
const socket = new WebSocket('ws://localhost:8080/ws');

socket.onopen = () => {
    console.log('WebSocket connection established');
};

socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};

socket.onclose = () => {
    console.log('WebSocket connection closed');
};

// Handle incoming messages from the WebSocket
socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    updateCanvas(data);
};

canvas.addEventListener('click', (event) => {
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    console.log('Canvas clicked at:', x, y);
    // Optionally, send this data to the backend
    // Example: Sending a message to the server
    socket.send(JSON.stringify({ type: 'click', x: x, y: y }));
});

function updateCanvas(data) {
    console.log(data)
    // Clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    // Draw simulation data
    fillColor = 0xff2244;
    data.forEach((entity) => {
        ctx.beginPath(); // Start a new path
        ctx.arc(entity.X, entity.Y, entity.Width, 0, Math.PI * 2); // Create a circle
        ctx.fillStyle = fillColor; // Set the fill color
        ctx.fill(); // Fill the circle
    });
}
