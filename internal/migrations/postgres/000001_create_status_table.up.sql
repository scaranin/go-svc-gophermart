CREATE TABLE status_order_info (
    status VARCHAR(20) PRIMARY KEY,
    date_start DATE,
    date_end DATE,
    is_active NUMERIC NOT NULL,
    is_final NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

insert 
  into status_order_info
values 
     ( 'NEW'
     , '0001-01-01'
     , '9999-12-31'
     , 1
     , 0
     , CURRENT_TIMESTAMP
     , CURRENT_TIMESTAMP
     ) ;

insert 
  into status_order_info
values 
     ( 'INVALID'
     , '0001-01-01'
     , '9999-12-31'
     , 1
     , 1
     , CURRENT_TIMESTAMP
     , CURRENT_TIMESTAMP
     ) ;

insert 
  into status_order_info
values 
     ( 'PROCESSING'
     , '0001-01-01'
     , '9999-12-31'
     , 1
     , 0
     , CURRENT_TIMESTAMP
     , CURRENT_TIMESTAMP
     ) ;

insert 
  into status_order_info
values 
     ( 'PROCESSED'
     , '0001-01-01'
     , '9999-12-31'
     , 1
     , 1
     , CURRENT_TIMESTAMP
     , CURRENT_TIMESTAMP
     ) ;
