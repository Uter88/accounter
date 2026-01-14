CREATE TABLE IF NOT EXISTS public.users (
    id BIGSERIAL NOT NULL,
    login CHARACTER VARYING(50) NOT NULL,
    password CHARACTER VARYING(50) NOT NULL,
    name CHARACTER VARYING(50) NOT NULL,
    surname CHARACTER VARYING(50) NOT NULL,
    patronymic CHARACTER VARYING(50) NOT NULL,
    price_per_hour NUMERIC(12, 2) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO public.users (id, login,"password",name,surname,patronymic,price_per_hour) VALUES
	(1, 'uter88@gmail.com','28031988','Евгений','Лихачев','Сергеевич',2366.29),
	(2, 'marynakoly@gmail.com','47873','Татьяна','Романченко','Сергеевна',1351.84) 
ON CONFLICT DO NOTHING;
