
# Booking-App Project with Beego Framework


The Booking App is designed to fetch and store hotel and rental property data from the Booking.com API into a local PostgreSQL database. It implements a backend API using the Beego framework to serve property data and provides a frontend for displaying property listings and details.

## Goals

1. **Database**: 
   - Fetch data from the Booking.com API.
   - Store data in a local PostgreSQL database with the following tables:
     - **Location Table**: Contains location names and their associated hotels.
     - **Rental Property Table**: Stores basic information about properties (name, type, bedrooms, bathrooms, amenities).
     - **Property Details Table**: Stores images, descriptions, and other relevant information about properties.

2. **API**:
   - Design a backend API with the following endpoints:
     - `/fetch_locations`: Fetch and store locations and associated hotels.
     - `/fetch_stays_data`: Fetch filtered stays data based on criteria.
     - `/fetch-hotel-details`: Fetch detailed information for a specific hotel.
     - `/fetch-hotel-images-and-description`: Fetch images and descriptions for a hotel.

3. **Frontend**:
   - Design two pages:
     - **Property Listing**: Displays all properties for a selected location with a search box for location search, a breadcrumb for navigation, and a tiles section for property information.
     - **Property Details**: Displays detailed information for a specific property.

## Prerequisites
Before you begin, ensure you have met the following requirements:

- **Go (Golang)**: Version 1.18 or higher.
- **Beego Framework**: Installed globally.
- **Git**: For cloning the repository.

## Installation and Setup

### Step 1: Install Go
   - Download and install Go from [https://go.dev/dl/](https://go.dev/dl/).
   - Verify installation:
     ```bash
     go version
     ```
### Step 2: Clone the repository
To get started, clone the project repository:
1. Navigate to your Go src `/go/src` directory to ensure the project is placed in the correct directory for your Go workspace.
2. Clone the repository:
```bash
git clone https://github.com/siddiqua14/booking-app.git
cd booking-app
```


### Step 3: Install Beego Framework
Beego is the framework used for this project, and Bee CLI is a development tool. If `Beego` is not installed in your workspace,
- Install them by running:
```bash
go get github.com/beego/beego/v2@latest
```
Ensure your `GOPATH` is set up correctly:
#### For Linux:
Add the following lines to your `~/.bashrc` or `~/.zshrc`:
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```
Reload your shell:
```bash
source ~/.bashrc
```
#### For Windows:
Set `GOPATH` and add it to your system `PATH`:
1. Open **Environment Variables** in System Properties.
2. Add a new variable:
   - Variable Name: `GOPATH`
   - Variable Value: `C:\Users\<YourUsername>\go` or in which workspace you like, setup path of that workspace
3. Edit the `Path` variable and add: `%GOPATH%\bin`
Verify your setup:
```bash
echo $GOPATH   # For Linux/MacOS
echo $env:GOPATH   # For Windows PowerShell
echo $env:Path 
```
##### To verify installation:
```bash
bee version
```


### Step 4: Docker
Set up the PostgreSQL database using Docker Compose:
- Ensure Docker and Docker Compose are installed on your machine.
- Run the following command to start the database:
```bash
docker-compose up 
```
- This command will create and start the PostgreSQL database container defined in docker-compose.yml.
- Access and make connection with the Databse 
```bash
docker exec -it postgres_booking psql -U postgres
\c booking_db
```
### Step 5: Install Dependencies

After installing the required Go modules and Beego dependencies:

```bash
go mod tidy
```
This will automatically resolve any missing dependencies and update your go.sum file with the required entries.
Run the application:
```bash
bee run 
```
Access the application at `http://localhost:8080`.
## API Endpoints

### Data Fetching Endpoints

- **Fetch Locations**: 
  - **Route**: `/fetch_locations`
  - **Controller**: `LocationController`
  - **Method**: `GET`
  - **Description**: Fetch and store locations and associated hotels.

- **Fetch Stays Data**:
  - **Route**: `/fetch_stays_data`
  - **Controller**: `LocationController`
  - **Method**: `GET`
  - **Description**: Fetch filtered stays data based on specified criteria.

- **Fetch Hotel Details**:
  - **Route**: `/fetch-hotel-details`
  - **Controller**: `LocationController`
  - **Method**: `GET`
  - **Description**: Fetch detailed information for a specific hotel.

- **Fetch Hotel Images and Descriptions**:
  - **Route**: `/fetch-hotel-images-and-description`
  - **Controller**: `ImageDescriptionController`
  - **Method**: `GET`
  - **Description**: Fetch images and descriptions for a hotel.

### Frontend Display Endpoints

- **Property Listing**:
  - **Route**: `/v1/property/list`
  - **Controller**: `PropertyController`
  - **Method**: `GET`
  - **Description**: Provides a list of properties for a specific location.

- **Property Details**:
  - **Route**: `/v1/property/details`
  - **Controller**: `PropertyController`
  - **Method**: `GET`
  - **Description**: Provides detailed information for a specific property.