### Description

The following code allows you to open a local server which will allow you to convert a sentence of your choice into an ascii-art in one of the banners available.

### Authors

The programme has been created by Romain Maillard and Nathan Paccoud.

### How to use it

To use the program, you must first open the file in a terminal and type "docker run -p 80:4000 --name ascii-art-app --rm ascii-art-web:1.0" .
This will start a local server you can access by typing "localhost:4000" on any browser. This will load a web page where you will be able to write a sentence of your choice. Note that in order to prevent server crashes, a limit on how many characters you may type has been added (1000 characters). Next, you can choose the banner you want to use by clicking on them. After clicking on the "envoyer" button, your phrase will be typed in the banner of your choice in the red square beneath.

### Algorithm details

The program is using a modified version of a program we created in the scope of another exercise. This program is based on two functions : we first transform each rune of the sentence into a slice of strings containing each of the eight lines of the letter in the banner selected. Since we chose to use the same model for the transformed sentenced, we created a function capable of adding a slice of rune to the end of another. But instead of using a custom sentence to print the result, we convert it back to a string to easily display it on the website.