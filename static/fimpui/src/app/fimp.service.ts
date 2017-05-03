import { Injectable } from '@angular/core';
import { Observable,Subject } from 'rxjs/Rx';
import {
  MqttMessage,
  MqttModule,
  MqttService
} from 'angular2-mqtt';
import { FimpMessage, NewFimpMessageFromString } from "app/fimp/Message";



@Injectable()
export class FimpService{
  private messages:FimpMessage[]=[];
  public observable: Observable<MqttMessage> = null;
  constructor(public mqtt: MqttService) {
     mqtt.onConnect.subscribe((message: any) => {
          console.log("FimService onConnect");
         // this.observable = null;
     }); 
     this.subscribeToAll("pt:j1/#");
  }
  public subscribeToAll(topic: string):Observable<MqttMessage>{
    this.observable = this.mqtt.observe(topic);
    this.observable.subscribe((msg) => {
      console.log("New message from topic :"+msg.topic+" message :"+msg.payload)
      this.saveMessage(msg);
    });
    return this.observable
  }
  public getGlobalObservable():Observable<MqttMessage>{
    if (this.observable == null){
      this.subscribeToAll("pt:j1/#");
    }
    return this.observable;
  }
  public subscribe(topic: string):Observable<MqttMessage>{
    return this.mqtt.observe(topic);
  }
  public publish(topic: string, message: string) {
    this.mqtt.publish(topic, message, {qos: 1}).subscribe((err)=>{
      console.log(err);
    });
  }
 private saveMessage(msg:MqttMessage){
      console.log("Saving new message to log")
      let fimpMsg  = NewFimpMessageFromString(msg.payload.toString());
      fimpMsg.topic = msg.topic;
      fimpMsg.raw = msg.payload.toString();
      fimpMsg.localTs =  Date.now();
      this.messages.push(fimpMsg);
 }

 public getMessagLog():FimpMessage[]{
   return this.messages
 }
}


@Injectable()
export class WsService {

    private actionUrl: string;
    private websocket: any;
    private receivedMsg: any;
    private observable:Observable<any>;
    
    constructor(){
      console.log("Fimp service constructor")
      this.connect();
      this.websocket = new WebSocket("ws://echo.websocket.org/"); //dummy echo websocket service
      this.websocket.onopen =  (evt) => {
          
          this.websocket.send("Hello World");
      };
    }

    public sendMessage(text:string){
      this.websocket.send(text);
    }

    public GetInstance(): Observable<any> {
      return this.observable
    }

    public connect(){
     this.observable = Observable.create(observer=>{
          this.websocket.onmessage = (evt) => { 
              observer.next(evt);
          };
      })
      .map(res=>"From WS: " + res.data)
      .share();
      // var subject = new Subject();
      // this.observable = source.multicast(subject);
      // this.observable.connect();
      
    }
}