WhatsOn is a simple Go command line program to interact with the Claude API (using go-anthropic) to return a number of TV programme viewing suggestions given a statement of desire.
The purpose of the program is to allow me to experiment with Claude embedding, its use of tools, and its XML input and output techniques.

The program uses a local database of TV schedules which Claude requests access to through its 'tool' mechanism.
To avoid having to recreate the database daily, dates in the database are automatically shifted to appear as if they start from the current date.

*Usage*
```sh
./whatson [-d] [-v] Statement of desire ...
```

The -d flag causes the program to follow the list of suggestions with a description of each program. Clause, being an AI, can of course make changes to those descriptions and this can be a way in which you can gain some insight into why it has made a suggestion.

For example:
```
./whatson I would like watch programmes about archaeology that are on in the next couple of weeks.

- Sun 2024-07-07 at 20:00: Digging for Britain on BBC4
- Sun 2024-07-14 at 20:00: Digging for Britain on BBC4
```

An explanation of Claude's thinking can be viewed by providing the -v flag.

```
./whatson -v I would like watch programmes about archaeology that are on in the next couple of weeks.

To provide appropriate TV programme suggestions about archaeology for the next couple of weeks, I'll need to use the get_tv_programmes function. Let's gather the necessary information:

1. The current date is 2024-06-30 (Sunday).
2. The user wants programmes about archaeology in the next couple of weeks.

I'll set the date range from today (2024-06-30) to two weeks later (2024-07-14) to cover the requested period.

Using the tool date range: 2024-06-30 to 2024-07-14
- Sun 2024-07-07 at 20:00: Digging for Britain on BBC4
- Sun 2024-07-14 at 20:00: Digging for Britain on BBC4
```

A description of each programme can be obtained with the -d flag.

```
./whatson "Show me suggestions for programmes that are related to animals and are on in the next couple of weeks."
- Sun 2024-07-02 at 18:00: Countryfile on BBC1
- Sun 2024-07-02 at 20:00: The Great British Bake Off on Channel 4
- Fri 2024-07-05 at 20:00: The Secret Life of the Zoo on BBC4
- Mon 2024-07-08 at 22:30: The Sky at Night on BBC4
- Sun 2024-07-14 at 20:00: Digging for Britain on BBC4
```

Asking Claude to explain itself through the descriptions can be insightful.
```
./whatson -d "Show me suggestions for programmes that are related to animals and are on in the next couple of weeks.  Update the programme description to show how it is related to animals."
- Sun 2024-07-02 at 18:00: Countryfile on BBC1
- Sun 2024-07-02 at 20:00: The Great British Bake Off on Channel 4
- Fri 2024-07-05 at 20:00: The Secret Life of the Zoo on BBC4
- Mon 2024-07-08 at 22:30: The Sky at Night on BBC4
- Sun 2024-07-14 at 20:00: Digging for Britain on BBC4

Descriptions:
- Countryfile : Exploring rural issues and celebrating the beauty of the British countryside, with segments featuring farm animals
- The Great British Bake Off : Amateur bakers compete in a series of challenging baking tasks using ingredients like eggs and dairy from animals.
- The Secret Life of the Zoo : Behind-the-scenes look at the lives of various animals and their keepers at Chester Zoo, providing intimate insights into animal behavior, care, and conservation efforts.
- The Sky at Night : While primarily about astronomy, this episode may explore the impact of space exploration on animal research, such as studying the effects of microgravity on various species or discussing animals used in early space missions.
- Digging for Britain : This archaeology program often uncovers animal remains and artifacts related to historical human-animal interactions, providing insights into ancient fauna and our evolving relationships with animals throughout British history.
```