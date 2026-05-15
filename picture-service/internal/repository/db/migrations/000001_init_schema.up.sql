CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS scenes (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    scene_index INT NOT NULL,
    text TEXT NOT NULL,
    image_url TEXT NOT NULL
);