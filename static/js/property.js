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

document.addEventListener('DOMContentLoaded', () => {
    // Get the property ID from the URL query string
    const urlParams = new URLSearchParams(window.location.search);
    const propertyID = urlParams.get('id');
    
    if (propertyID) {
        fetchPropertyDetails(propertyID);
    } else {
        console.error('Property ID is missing in the URL');
    }
});

function fetchPropertyDetails(propertyID) {
    fetch(`/v1/property/details?id=${encodeURIComponent(propertyID)}`)
        .then(response => response.json())
        .then(data => {
            if (data) {
                displayPropertyDetails(data);
            } else {
                console.error('No data received or property not found');
            }
        })
        .catch(error => console.error('Error fetching property details:', error));
}

function displayPropertyDetails(property) {
    const propertyDetailsContainer = document.querySelector('.property-details');
    propertyDetailsContainer.innerHTML = `
        <h1>${property.HotelName}</h1>
        <p><strong>Location:</strong> ${property.Location}</p>
        <p><strong>Price:</strong> ${property.Price}</p>
        <p><strong>Type:</strong> ${property.PropertyType}</p>
        <p><strong>Guests:</strong> ${property.Guests}</p>
        <p><strong>Rating:</strong> ${property.Rating} (${property.ReviewCount} Reviews)</p>
        <p><strong>Beds:</strong> ${property.NumBeds}</p>
        <p><strong>Bedrooms:</strong> ${property.NumBedR}</p>
        <p><strong>Bathrooms:</strong> ${property.NumBaths}</p>
        <p><strong>Bedrooms:</strong> ${property.Bedroom}</p>
        <p><strong>Amenities:</strong> ${property.Amenities.join(', ')}</p>
        <div class="images">
            ${property.Images.map(image => `<img src="${image}" alt="Property Image" class="property-image">`).join('')}
        </div>
        <p>${property.Description}</p>
    `;
}