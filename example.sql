-- Get owner of %(parcel_id)s
SELECT owner.fullname FROM owner JOIN parcel ON (owner.parcel_id = parcel.id)
    WHERE parcel.parcelid = %(parcel_id)s
    ORDER BY owner.creationts
    DESC limit 1;
