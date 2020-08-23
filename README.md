## 목적  
DSC에서 프로젝트를 진행함에 있어. Discord를 자주 사용하게되었다. (시국)   
그래서 DSC Discord 채널에서 사용할 수 있는 DiscordBot을 만드려고한다.
프로젝트 저장소는 따로 있지만 마구 커밋하기에는 조금 쑥스러우니 개인 저장소에서 막 코드작성하다가
나중에 있어보이게 뭉텅이로 보낼 예정입니다.


## 진행사항
2020.08.19 - postgresql을 이용한다. DiscordGo를 이용한다. 채널에서 챗봇에서 "item [정보1] [정보2]"를 입력하면 DB에 정보가 insert 된다. 다시 명령어를 사용하면 터진다. primary 키로 채널id를 넣어뒀기 때문이다. 다시 !item입력시 update되도록 바꿔야 한다.  

2020.08.23 - 이전에 update 문제가 있었다. 해당 채널이 데이터베이스에 등록되어있는지를 'activate' 항목으로 확인할 수 있고, 'activate'가 true 일 경우 새로 insert하는게 아닌 update문을 실행시키는 형식을 취했다. 그런데 postgresql에 update에 대해 내가 이해가 부족하 것인지. query문이 작동은 하는 것 같은데 데이터베이스 update가 되지 않고 있다. 그래서 오늘 그 부분으 고치려고 하는 중이다.(ing)


## 참고블로그
[golang postsql querry 문짜기](https://brownbears.tistory.com/186). 
[psql server error](https://velog.io/@kim-macbook/postgresql-error-1)
