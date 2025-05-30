create index idx_users_points on users(points desc);
create index idx_users_city_points on users(city, points desc);
create index idx_friends_user_id on friends(user_id);
create index idx_friends_friend_id on friends(friend_id);
//3.1
create materialized view top100_users as
select id, username, city, points
from users
order by points desc
    limit 100;
