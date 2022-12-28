# Twitch Trivia/Scramble Solver

Twitch Trivia/Scramble Solver will automatically attempt to solve trivia questions and scramble questions asked by Amazefulbot in Twitch chatrooms. If it cannot solve the question, it will learn the question and the answer to said question.

## Instructions

Download [`twitch-trivia-unscrambler.zip`](https://github.com/ActuallyGiggles/twitch-trivia-scramble-solver/releases/tag/1.0.0) and unzip it. Now you can launch `twitch-trivia-unscrambler.exe` inside of that folder and the program will run you through the first time setup.

## Additional Information

1. You can specify whether to answer trivia, scramble, or both.
2. You can specify which channels to answer in.
3. You can specify how long to wait before answering (to seem more human).
4. You can specify the percentage of questions that should be ignored (to seem more human).
5. For trivia, you can specify the percentage of answering partially first (to seem more human).

In addition to downloading the program, you must also download two JSON files. The first file contains a list of many trivia questions and answers. This is the largest file. There is no guarantee that every question is found in it, but I tried my best to gather what I could from existing questions and over 200,000+ jeopardy questions. The second file contains a list of common words that might be scrambled.

If you would like to add channels to monitor, edit the config JSON file.
