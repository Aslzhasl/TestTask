select * from users order by points desc;
//2.1
select
    id,
    username,
    city,
    points,
    rank() over (order by  points desc) as global_rank
from users
order by points desc;
//2.2
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
