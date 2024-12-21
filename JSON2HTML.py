from json2html import *
import json

jsonfilename = "app_bsky_feed_post.json" 

htmlfilename = "repo.html" 

    
with open(jsonfilename, "r", encoding='utf-8') as jasonfile:
    JasonData = json.load(jasonfile)

# unneeded using with  
# jasonfile.close

  
# print(json2html.convert(json = JasonData))

with open(htmlfilename, "w", encoding='utf-8') as htmlfile:
    htmlfile.write(json2html.convert(json = JasonData))

# unneeded using with  
#htmlfile.close
