-- setup.sql

-- REMOVE ALL DATA FROM THE TABLES ON THE EACH SCRIPT RUN, IF TABLES EXISTS
DO $$
BEGIN
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'jobs') THEN
        TRUNCATE TABLE jobs RESTART IDENTITY CASCADE;
    END IF;
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'technologies') THEN
        TRUNCATE TABLE technologies RESTART IDENTITY CASCADE;
    END IF;
END $$;

-- Create technologies table
CREATE TABLE IF NOT EXISTS technologies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    djinni_keyword VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title TEXT,
    company TEXT,
    work_format TEXT,
    location TEXT,
    company_type TEXT,
    experience_years NUMERIC,
    english_level TEXT,
    technology INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (technology) REFERENCES technologies(id) ON DELETE SET NULL
);

INSERT INTO technologies (name, djinni_keyword) VALUES
  ('Angular', 'Angular'),
  ('React.js', 'React.js'),
  ('Svelte', 'Svelte'),
  ('Vue.js', 'Vue.js'),
  ('Java', 'Java'),
  ('C# / .NET', '.NET'),
  ('Python', 'Python'),
  ('WordPress', 'WordPress'),
  ('Yii', 'Yii'),
  ('Drupal', 'Drupal'),
  ('Laravel', 'Laravel'),
  ('Magento', 'Magento'),
  ('Symfony', 'Symfony'),
  ('Node.js', 'Node.js'),
  ('iOS', 'iOS'),
  ('Android', 'Android'),
  ('React Native', 'React%20Native'),
  ('C', 'C%20Lang'),
  ('Embedded', 'Embedded'),
  ('C++', 'CPP'),
  ('Flutter', 'Flutter'),
  ('Golang', 'Golang'),
  ('Ruby', 'Ruby'),
  ('Scala', 'Scala'),
  ('Salesforce', 'Salesforce'),
  ('Rust', 'Rust'),
  ('Elixir', 'Elixir'),
  ('Kotlin', 'Kotlin'),
  ('MS Dynamics / Business Central', 'MS%20Dynamics'),
  ('Odoo', 'Odoo'),
  ('SAP', 'SAP')
ON CONFLICT (name) DO NOTHING;
