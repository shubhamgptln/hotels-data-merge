# hotels-data-merge

## **Steps To Run**

1. Run the following command :-

        make run-api
2. Make a curl request :-
        
         `curl -v 'http://localhost:8090/hotels-data/filter' --data '{"hotel_id": ["iJhz","f8c9","SjyX"],"destination_id": {fill}}'

## **Merge Strategy**
1. Name is the longest string among all the supplier entries.
2. Description is the longest string among all the supplier entries.
3. Booking Conditions is a list of unique Booking condition statements.
4. Location is chosen based on which entry has most non-null(default) values in Location struct
5. Images is a list of unique image urls for each of site,room and amenities category. The image description is chosen as the longest string for the same image url.
6. Amenities are merged based on unique values.Uppercase and lowercase characters are considered same. If an amenity exists in room category as well as general category we remove that entry from general category as that amenity is best suited for room category.
7. Choice of merge strategies are made configurable via spec file.
8. We have a model layer where we define generic merging strategies irrespective of actual fields.

## **Sanitization**

I have used a simple typecast for fields which didn't meet the expected model layer schema. The dirty fields were Lat and Lng fields from one of the suppliers. This logic can be applied to other fields as well if we have no control over the data types and values of each fields. For now I just assumed that Lat and Lng fields are the only dirty fields.

## **Procurement**

I have created a HTTP client for each of the suppliers and called them parallely using goroutines to fetch data quickly. This data is then combined and stored in map with key as hotelID and value as array containing json data from all the suppliers.

# **Containerization**

There is a DockerFile to containerize the application. Steps to run app container are as follows:-
1. To build a container 

   `docker build -t hotel-merge-app -f ./deployments/Dockerfile/ .` 
