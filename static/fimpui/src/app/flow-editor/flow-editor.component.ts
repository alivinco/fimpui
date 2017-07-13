import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Http, Response,URLSearchParams }  from '@angular/http';

@Component({
  selector: 'app-flow-editor',
  templateUrl: './flow-editor.component.html',
  styleUrls: ['./flow-editor.component.css']
})
export class FlowEditorComponent implements OnInit {
  flow :Flow;
  constructor(private route: ActivatedRoute,private http : Http) {
    this.flow = new Flow();
   }

  ngOnInit() {
    let id  = this.route.snapshot.params['id'];
    this.loadFlow(id);
  }

  loadFlow(id:string) {
     this.http
      .get('http://localhost:8081/fimp/flow/definition/'+id)
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         this.flow = result;
         console.dir(this.flow)
         console.log(this.flow.Name)
      });
  }
  saveFlow() {
    console.dir(this.flow)
  } 

}
export class Flow {
    Id :string ;
    Name : string ;
    Description : string ;
    Nodes : any[] ;
}