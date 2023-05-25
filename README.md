# Saul

Saul is a prototype application that uses gpt3.5 to create lesson plans for a given grade level, student population, and item descriptor (topic).

There are currently 2 "modes" in Saul. The `free entry` mode allows users to enter whatever combination of grade level, student population, and item descriptor they choose into the app. `guided` mode connects to a demo [Firestore](https://cloud.google.com/firestore) database that contains sample data for 3 hypothetical schools. In this mode, Saul will guide users through a process of choosing a school and choosing a test, upon which it will recommend the three "highest leverage" areas for instruction and allow users to generate lessons aligning with those areas.

Saul is currently a work in progress and is subject to change.