-- LoadImage Загружаем изображение в бд
-- name: LoadImage :one
INSERT INTO images (name, content, id_user)
VALUES (@name::text, @content::bytea, @id_user::integer)
RETURNING *;

-- DeleteImage Удаляем картинку
-- name: DeleteImage :one
DELETE FROM images
WHERE id_image = @id_image::integer
RETURNING *;

-- GetImage Возвращаем картинку
-- name: GetImage :one
SELECT * FROM images
WHERE id_image = @id_image::integer;

-- RenameImage Переименовываем картинку
-- name: RenameImage :one
UPDATE images
SET name = @name::text
WHERE id_image = @id_image::integer
RETURNING *;

-- EditImage Заменяем картинку на новую
-- name: EditImage :one
UPDATE images
SET content = @content::bytea
WHERE id_image = @id_image::integer
RETURNING *;
