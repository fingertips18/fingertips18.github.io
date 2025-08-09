import { CMES, KTMSCES, MNHS, USM } from './assets';
import { EDUCATIONTYPE } from './enums';

export const EDUCATIONS = [
  {
    source: 'https://www.usm.edu.ph',
    name: 'University of Southern Mindanao',
    logo: USM,
    department: 'College of Engineering and Information Technology',
    date: 'August 2019 - June 2023',
    honor: 'Cum Laude',
    desc: "I graduated from the University of Southern Mindanao (USM) with a Bachelor's degree in Computer Science, achieving Cum Laude with a GPA of 1.67. During my time at USM, I actively participated in the in-house review, representing my thesis study. I was also a member of the Philippine Society of Information Technology (PSIT), which enriched my academic experience. My coursework included Data Structures and Algorithms, Time Complexity, Software Engineering, Networking, and Artificial Intelligence, among other advanced topics.",
    study: {
      title:
        'Luminous: A heart rate-based horror adventure game using A* pathfinding algorithm',
      desc: "Luminous is a story-based horror-adventure game that I created in Unity for my undergraduate thesis. The study mainly focused on AI (artificial intelligence) and provided a unique mechanism for tracking the player's location based on their heart rate. The objectives of this study were to use the player’s heart rate as the heuristic value in the A* algorithm, implement an enemy-tracking mechanic based on the heart rate-based heuristics, and determine its accuracy against the default A* algorithm.",
      stack: [
        'Unity',
        'A*',
        'C#',
        'HypeRate',
        'Photoshop',
        'Blender',
        'Audacity',
      ],
      demo: 'https://www.youtube.com/watch?v=7zYUk5x-B40',
    },
    projects: [
      {
        title: 'Mastivity',
        desc: 'As per my OJT requirement, we were tasked with creating a system dedicated to our assigned department. I was assigned to graduate school; thus, I created a system that boosts masters productivity and will help them with their daily endeavors.',
        stack: [
          'Bootstrap',
          'Angular',
          '.Net Core',
          'Entity Framework',
          'Swagger API',
          'MSSQL',
          'Azure',
          'Netlify',
        ],
        demo: 'https://www.youtube.com/watch?v=OUnh-eysJrM',
      },
      {
        title: 'Document Request System',
        desc: 'This was a system that I made for the HR department at USM for generating documents based on user requests by filling out an online form provided by the system.',
        stack: [
          'Bootstrap',
          'Angular',
          '.Net Core',
          'Entity Framework',
          'Swagger API',
          'MSSQL',
          'Azure',
          'Netlify',
        ],
        demo: 'https://www.youtube.com/watch?v=jkJ1Z9-yHYU',
      },
      {
        title: 'Faculty Competency System',
        desc: 'This was a system that I made for the HRDMO to assess the competency level of the faculty members.',
        stack: [
          'Bootstrap',
          '.Net Blazor',
          'Entity Framework',
          'Swagger API',
          'MSSQL',
          'Azure',
          'Netlify',
        ],
      },
    ],
    type: EDUCATIONTYPE.college,
  },
  {
    source: 'https://www.facebook.com/MatanaoNHS',
    name: 'Matanao National High School',
    logo: MNHS,
    department:
      'Senior - Information and Communication Technology (ICT) Strand',
    date: 'June 2017 - April 2019',
    honor: 'With High Honors',
    desc: 'I graduated from Matanao National High School (MNHS) - Senior High with the distinction of With High Honors. This achievement reflects my dedication and commitment to academic excellence throughout my senior high school years.',
    type: EDUCATIONTYPE.seniorHigh,
  },
  {
    source: 'https://www.facebook.com/MatanaoNHS',
    name: 'Matanao National High School',
    logo: MNHS,
    department: 'Junior - Science, Technology, Engineering and Mathematics',
    date: 'August 2013 - June 2017',
    honor: 'With Honors',
    desc: 'I completed my junior high school education at Matanao National High School (MNHS), graduating with the distinction of With Honors. This recognition highlights my consistent academic performance and dedication during those formative years.',
    type: EDUCATIONTYPE.juniorHigh,
  },
  {
    source:
      'https://www.facebook.com/p/DepEd-Tayo-Youth-Formation-Ceboza-Elementary-School-100079755368493/?_rdr',
    name: 'Ceboza Matanao Elementary School',
    logo: CMES,
    sub: {
      name: 'Kapitan Tomas Monteverde Sr. Central Elementary School',
      desc: 'Grade 1 - 5',
      logo: KTMSCES,
    },
    date: 'June 2007 - March 2013',
    honor: 'Valedictorian',
    desc: 'I began my elementary education at Kapital Tomas Monteverde Sr. Central Elementary School, and later transferred to Ceboza Matanao Elementary School in Grade 5. I graduated from Ceboza Matanao Elementary School as the class Valedictorian, an honor that reflects my commitment to academic excellence from an early age.',
    type: EDUCATIONTYPE.elementary,
  },
];
