function searchLocation() {
    const location = document.getElementById('location-search').value;
    const breadcrumb = document.getElementById('breadcrumb-value');
    breadcrumb.textContent = location; // Update the breadcrumb

    // Fetch property data based on location
    fetch(`/property/list?location=${location}`)
        .then(response => response.text())
        .then(html => {
            document.body.innerHTML = html; // Update the page with new HTML from the server
        })
        .catch(error => console.error('Error fetching property data:', error));
}

// Trigger search on page load
window.onload = () => {
    searchLocation();
};

document.addEventListener('DOMContentLoaded', function() {
    // Add any JavaScript interactions here
    const searchForm = document.querySelector('.search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', function(e) {
            e.preventDefault();
            // Handle search form submission
        });
    }
});
