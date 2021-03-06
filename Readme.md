#Approach

* Picks the longest found string in the dictionary and converts the same to title case (except the first word it finds)
* For instance, in case of "onetwofour" , the following happens
    1. startIndex points to the first character and endIndex to the last character.
    2. endIndex is decremented until the string between startIndex and endIndex is found in the dictionary.
    3. startIndex then goes to endIndex+1 and steps 1 to 3 are repeated until startIndex goes beyond the last character in the string.

* TitleCase of words found in the dictionary are cached for a certain duration to reduce the number of hits to the dictionary API.
* The app has been containerized and hosted on Azure Container Registry which is where the app service is pointing to (dockerFile is included).
