**simpleSurance**

Project content:
project contains 5 folders, one main file and one Readme
DB           - Folder where I save text files for backup
FileManager  - manager: is responsible for communication with files.
               manager_test: for covering some test cases.
MovingWindow - this package is responsible for maintaining correct list
               of time records in runtime.
               window_test: for covering some test cases.
               
Pages        - Simple one html template for showing result
Server       - Server, which creates and holds file manager and window.

Usage:
main.go - Just creates server objects with address and pattern and starts it

BackUp logic:
*For safe backup I use 2 files,
*for each file I have defined time, when I modified it last time. 
*from the beginning I start writing in file1
*Change active file when another file becomes older than 1 minute, this means
 I can clear it completely, and start writing in it
*I always write either in file1 or file2, this means I do not have intersections,
 records are sorted, to get sorted records list I just need to concatenate them,
 they maybe contain some old records but it is ok, we can filter it in window.
 on backup I continue writing in file which modified time is bigger.
 
 I decide to not to handle requests in different go routines, because here we have a 
 bottleneck-writing in file, this means only one routine can write in it at the same time,
 so I decide to not to use them. I had an idea to use pool of free files, and make backup 
 from several files but I decide to not make it complex.