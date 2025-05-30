select id, username, points
from users
order by points desc
    limit 100;
--4.1
--Тут всё просто: чтобы получить топ пользователей или топ по городу — сортируем по points,
--лимитируем количество, и всё. Для топа среди друзей просто соединяем таблицы и сортируем.
--С хорошими индексами всё летает даже на миллионах строк
select id, username, points
from users
where city = 'Дубай'
order by points desc
    limit 100;
--4.2
select u.id, u.username, u.points
from friends f
         join users u on u.id = f.friend_id
where f.user_id = (select id from users where username = 'user2284483')
order by u.points desc
    limit 100;
