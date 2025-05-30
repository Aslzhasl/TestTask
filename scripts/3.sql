create index idx_users_points on users(points desc);
create index idx_users_city_points on users(city, points desc);
create index idx_friends_user_id on friends(user_id);
create index idx_friends_friend_id on friends(friend_id);
--3.1
--Оптимизация здесь — это индексы. Индекс по points ускоряет сортировку и поиск топов.
--Индекс по (city, points) нужен, чтобы топы по городам не тормозили.
--В friends ставим индексы по user_id и friend_id, чтобы быстро искать друзей и делать join'ы


create materialized view top100_users as
select id, username, city, points
from users
order by points desc
    limit 100;
