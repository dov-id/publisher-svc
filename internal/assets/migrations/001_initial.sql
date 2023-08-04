-- +migrate Up

create type requests_status_enum as enum ('success', 'in progress', 'failed', 'pending');

create table if not exists requests (
    id uuid primary key,
    status requests_status_enum not null,
    error text not null
);

create table if not exists feedbacks (
    course text not null,
    content text not null,

    unique(course, content)
);

create index feedbacks_course_idx on feedbacks(course);

-- just mock
INSERT INTO feedbacks values
    ('0x51F32e07441f8FA2F61B5Bc0368917faaDca1c98', 'The course was an excellent introduction to Python. The instructor explained concepts clearly and provided practical examples to reinforce learning. I now feel confident to start writing my own Python scripts.');

INSERT INTO feedbacks values
    ('0x6Ac7bC6c7BeEC5693c08734729Ed15fbc23D98c7', 'The course was a great introduction to DevOps principles and practices. I learned about continuous integration, deployment, and monitoring. This knowledge is valuable for modern software development.');

INSERT INTO feedbacks values
    ('0x29005cD047DBBB29C2a7ed797EC7C374FF344F46', 'The course provided valuable insights into data visualization techniques. I learned how to create visually appealing charts and graphs to effectively communicate data insights.');

INSERT INTO feedbacks values
    ('0x6f317EEa554B931FC405433A27C9d37A8cF0c9c8', 'The course was well-structured and covered essential Linux commands and administration tasks. The instructor''s explanations were concise and easy to follow. I now feel comfortable working in a Linux environment.');

INSERT INTO feedbacks values
    ('0xf3eaf3A53E99Fa68E4d395843303901548a1e9D8', 'The course provided a broad overview of AI concepts. I learned about machine learning, natural language processing, and robotics. It sparked my interest in further exploring AI technologies.');

INSERT INTO feedbacks values
    ('0x6f317EEa554B931FC405433A27C9d37A8cF0c9c8', 'This course was an excellent introduction to React. The instructor explained React components and state management clearly. The real-world projects helped me improve my front-end development skills.');

INSERT INTO feedbacks values
    ('0x6f317EEa554B931FC405433A27C9d37A8cF0c9c8', 'I found the course very informative and practical. The instructor covered all aspects of SQL databases, and I now feel comfortable writing SQL queries and managing databases.');

INSERT INTO feedbacks values
    ('0x29005cD047DBBB29C2a7ed797EC7C374FF344F46', 'The course provided a solid foundation in networking concepts. The instructor explained complex topics with clarity. Now I understand how data flows over networks and how to troubleshoot common issues.');

INSERT INTO feedbacks values
    ('0x29005cD047DBBB29C2a7ed797EC7C374FF344F46', 'The course was a great refresher on data structures and algorithms. The instructor''s approach to problem-solving was commendable, and I feel more confident in coding efficient solutions.');

INSERT INTO feedbacks values
    ('0x6Ac7bC6c7BeEC5693c08734729Ed15fbc23D98c7', 'This course gave me a clear understanding of cloud computing concepts and platforms. The practical exercises on cloud services were beneficial in learning how to deploy applications in the cloud.');

INSERT INTO feedbacks values
    ('0x6Ac7bC6c7BeEC5693c08734729Ed15fbc23D98c7', 'The course was well-structured and covered everything needed to create Android apps. The hands-on projects helped me grasp the concepts effectively. Now I can develop basic Android applications.');

INSERT INTO feedbacks values
    ('0x51F32e07441f8FA2F61B5Bc0368917faaDca1c98', 'This course was eye-opening and provided insights into the world of cybersecurity. I learned various hacking techniques and how to secure computer systems. It was a challenging but rewarding experience.');

INSERT INTO feedbacks values
    ('0x51F32e07441f8FA2F61B5Bc0368917faaDca1c98', 'The course provided a solid foundation in machine learning concepts. The instructor explained complex algorithms in a way that was easy to understand. I enjoyed the practical projects which allowed me to apply my learning.');

INSERT INTO feedbacks values
    ('0x51F32e07441f8FA2F61B5Bc0368917faaDca1c98', 'The instructor was knowledgeable and engaging. I learned how to perform data analysis and visualization using R. The hands-on exercises were particularly helpful in reinforcing the concepts.');

INSERT INTO feedbacks values
    ('0x6f317EEa554B931FC405433A27C9d37A8cF0c9c8', 'This course was very comprehensive and covered all the essential aspects of web development. I gained valuable knowledge in HTML, CSS, and JavaScript, and can now create basic web pages with ease.');

-- +migrate Down

drop index if exists feedbacks_course_idx;
drop table if exists feedbacks;
drop table if exists requests;
drop type if exists requests_status_enum;