CREATE TABLE education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    main_school JSONB NOT NULL,
    school_periods JSONB,
    projects JSONB,
    level TEXT CHECK (level IN (
        'elementary',
        'junior-high-school',
        'senior-high-school',
        'college'
    )) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
