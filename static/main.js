const canvas = document.getElementById('simulationCanvas');
const ctx = canvas.getContext('2d');

canvas.addEventListener('click', (event) => {
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    console.log('Canvas clicked at:', x, y);
    // Optionally, send this data to the backend
});

function updateCanvas(data) {
    // Clear the canvas
    console.log(data);
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    // Draw simulation data
    ctx.fillStyle = 'red';

    data.forEach(x => {
        ctx.fillRect(x.X, x.Y, x.Width, x.Height);
    })
}

async function fetchSimulationData() {
    try {
        const response = await fetch('/api/simulate');
        const data = await response.json();
        updateCanvas(data);
    } catch (error) {
        console.error('Error fetching simulation data:', error);
    }
}

// Periodically fetch and update simulation data
setInterval(fetchSimulationData, 1000);
