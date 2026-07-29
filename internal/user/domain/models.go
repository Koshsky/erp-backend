package domain

const (
	ProjectDirector string = "ДП"
	ProjectManager  string = "РП"
	ProcessOwner    string = "ВП"
)

/*
ProjectDirector (ДП) просматривает, создает и удаляет проекты;
Редактирует: code, priority, start_date, end_date,
назначает РП для проекта (owner_id)
+просматривает все процессы и задачи

ProjectManager (РП) просматривает свои проекты, процессы
создает и удаляет процессы,
редактирует процессы: title, start_date, end_date, owner_id
+ просматривает свои процессы и задачи к ним

ProcessOwner (ВП) просматривает свои процессы, задачи, ресурсы (общие)
Создает и удаляет задачи,
Редактирует задачи: title, start_date, end_date
Назначает и изменяет ресурсы для задач (assignment)

итого
ПРОЕКТЫ
просмотр
if role == "ДП" OR project.owner_id == user_id
создание/изменение/удаление
if role == 'ДП'

ПРОЦЕССЫ
просмотр
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
создание/изменение/удаление
if project.owner_id == user_id

ТАСКИ
просмотр
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
создание/изменение/удаление
if process.owner_id == user_id
ASSIGNMENTS
просмотр
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
создание/изменение/удаление
if process.owner_id == user_id
RESOURCES
просмотр
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
создание/изменение/удаление // TODO: кто?
if process.owner_id == user_id
*/

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}
