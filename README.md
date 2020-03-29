## Welcome


Borges.ai is an experiment to build better Goodreads experiment.


## Next task

 - use metrics to report to the cloudwatch
 - add buttons to read book when browsing other peoples pages
 - settings: username
 - specify location and specify email in setting. Can we do autocomplete for location? 
 - signup via email
 - can I use HTTP2 to send some HTML for long requests? when I do post for sync?
 - improve author page https://borges.ai/a/neil_postman
 - save last visit date for user
 - https://bookshop.org/signup for the future. when I have 1000 users.
 - implement view counters fo user pages and user book pages. IP address. so we can do country later. maybe publish job and have geo location for IP
 - implement counters for visited users
 - implement twitter cards and share my reviews once I finish book on twitter - https://developer.twitter.com/en/docs/tweets/optimize-with-cards/guides/getting-started
 - add comma when rendering multiple authors on user page
 - search should be by author too. update placeholder and search for all the books by author in parallel
 - improve text editor! this one maybe https://github.com/sparksuite/simplemde-markdown-editor. or just code mirror. Can do fullscreen mode too. Need something that works on mobile.
 - integrated Google books API https://developers.google.com/books/docs/v1/using. Better search -> another parallel. Also more metadata and more numbers
    - lccn -  Library of Congress Control Number.
    - oclc - Online Computer Library Center number.
  
## Later

 - What is this https://www.worldcat.org/title/flatland-a-romance-of-many-dimensions/oclc/233582656? anything interesting there?
 - For some reason for some books (like Camus' Exile and the Kingdom) we can't fetch all needed data - like pages.
But it's there. We cannot get all editions from Goodreads API, but we can scrap https://www.goodreads.com/work/editions/1028050-l-exil-et-le-royaume?filter_by_format=&sort=original_date_published&utf8=✓.
 - privacy policy
 - terms
 - sitemap.xml per user. so @podviaznikov/sitemap.xml

## Tech

 - find how many percent has 1) twitter, 2) goodreads 3) both
 - how many people without goodreads created review
 - how many has goodreads_count > 0 but no reviews - some bug with permissions
 
## Promos

 - find more reading influencers. Substack, telegramm twitter - ask them to migrated to Borges. How much would they charge?
 - add on reddit about books and book clubs 
 
## Custom domains

There are limit of 25 rules. 
After 25 rules used create new ALB. We can do as many as needed.
 
## Promos

 - can I create telegram channel with discussion
 - can I create soundcloud channel to send my tweets/quotes as audio files?
 
## Ideas

 - how to mark if book is owned. Can we get that data from Goodreads initially? We need this later in order to show books that can be borrowed
 - we need to shop button. want to read if logged in on other users pages.
   

## Deploy

```
    cd frontend
    npm build
    cd ..    
    $(aws ecr get-login --no-include-email --region us-west-2 --profile borges)
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .
    docker build -t borges-web .
    docker tag borges-web:latest 411700958227.dkr.ecr.us-west-2.amazonaws.com/borges-web:latest
    docker push 411700958227.dkr.ecr.us-west-2.amazonaws.com/borges-web:latest
```
