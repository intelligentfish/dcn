create database if not exists dcn;
grant all on dcn.* to 'dev'@'%' identified by 'dev';
flush privileges;