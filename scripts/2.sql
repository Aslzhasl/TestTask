--Для получения всех пользователей по очкам просто сортируем таблицу по points дабы база быстрее стработал
select * from users order by points desc;
--2.1Чтобы узнать место каждого пользователя в общем рейтинге, используем оконную функцию rank() она сам все сдлеат
select
    id,
    username,
    city,
    points,
    rank() over (order by  points desc) as global_rank
from users
order by points desc;
--2.2Если нужно вытащить инфу по одному пользователю — типа его место глобально,
--в своём городе и среди друзей — делаем временную таблицу (через with), где сразу считаем все нужные ранги.
--Место среди друзей — просто считаем, сколько из его друзей набрали больше очков.
                                                                                                                                                                                                               Всё делается одним большим SQL-запросом, не надо ничего пихать в память приложения.
with ranked as (
    select
        id,
        username,
        city,
        points,
        rank() over (order by points desc) as global_rank,
        rank() over (partition by city order by points desc) as city_rank
    from users
)
select
    r.username,
    r.points,
    r.global_rank,
    r.city_rank,

    (select count(*) + 1
     from users u2
              join friends f on f.user_id = r.id
     where u2.id = f.friend_id and u2.points > r.points
    ) as friends_rank
from ranked r
where r.username = 'user123';

