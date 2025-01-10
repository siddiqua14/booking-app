document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.getElementById('search-button');
    searchButton.addEventListener('click', () => {
        const location = document.getElementById('locationSearch').value;
        fetchProperties(location);
    });
});

function fetchProperties(location = 'New York') {
    fetch(`/v1/property/list?location=${encodeURIComponent(location)}`)
        .then(response => response.json())
        .then(data => {
            if (data) {
                displayProperties(data);
            } else {
                console.error('No data received');
            }
        })
        .catch(error => console.error('Error fetching properties:', error));
}

function displayProperties(properties) {
    const propertyGrid = document.getElementById('propertyGrid');
    propertyGrid.innerHTML = '';

    properties.forEach(property => {
        const propertyCard = `
            <div class="bg-white rounded-lg shadow-md overflow-hidden">
                <div class="relative">
                    <img src="${property.HeroImage || '/static/images/placeholder.jpg'}" alt="${property.HotelName}" class="w-full h-64 object-cover">
                </div>
                <div class="p-4">
                    <div class="flex items-center mb-2">
                        <span class="text-blue-600 font-bold">${property.Rating.toFixed(1)}</span>
                        <span class="ml-2 text-gray-600">(${property.ReviewCount} Reviews)</span>
                    </div>
                    <h3 class="text-xl font-bold mb-2">
                        <a href="/property-details?id=${property.HotelID}" target="_blank" class="text-blue-900 hover:underline">
                            ${property.HotelName}
                        </a>
                    </h3>
                    <div class="text-gray-600 mb-2">
                        ${property.Location}
                    </div>
                    <div class="flex items-center justify-between mt-4">
                        <div class="text-lg font-bold">From ${property.Price}</div>
                        <button onclick="window.open('/property-details?id=${property.HotelID}', '_blank')" class="bg-emerald-500 text-white px-4 py-2 rounded-lg">
                            View Availability
                        </button>
                    </div>
                </div>
            </div>
        `;
        propertyGrid.innerHTML += propertyCard;
    });
}