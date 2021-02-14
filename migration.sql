create table if not exists realestateportal(
   id  int generated always as identity primary key ,
   creationts   timestamp default now(),
   url          text not null ,
   address      bytea[],
   municipality bytea[],
   owner        bytea[],
   changemail   bytea[]
);
comment on table realestateportal is 'Scraped data from the Allegheny County Real Estate Portal';

create table if not exists parcel(
    -- I believe we can use parcel id as the primary key. Things break if there are multiple parcels with the same id.
    -- However, a surrogate key ("id") is used as requested.
    id         int generated always as identity primary key ,
    parcelid   text unique,
    creationts timestamp default now()
);

create table if not exists owner(
    id                 int generated always as identity primary key ,
    creationts         timestamp default now(),
    parcel_id          int,
    fullname           text[] not null,
    bobsource_sourceid int
);

create table if not exists address (
   id                    int generated always as identity primary key ,
   parcel_id             int,
   fulladdress           text,
   bobsource_sourceid    int
);


alter table address
    add constraint address_parcel_id_fk
        foreign key (parcel_id)
            references parcel (id);

alter table address
    add constraint address_bobsource_sourceid_fk
        foreign key (bobsource_sourceid)
            references bobsource (sourceid);

alter table owner
    add constraint owner_bobsource_sourceid_fk
        foreign key (bobsource_sourceid)
            references bobsource (sourceid);

alter table owner
    add constraint owner_parcel_id_fk
        foreign key (parcel_id)
            references parcel (id);


-- alter table bobsource add primary key (sourceid);
-- INSERT INTO bobsource (sourceid, title, description, creator, muni_municode)
-- VALUES (20, 'SnapScrape', 'Scraped from Allegheny County Real Estate Portal', 100, 999);


