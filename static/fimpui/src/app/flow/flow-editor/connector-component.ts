import { Component, OnInit,Inject } from '@angular/core';

@Component({
    selector: 'conn-endpoint',
    template:`
    <svg style="position:absolute;left:0px;top:0px"
            width="800" height="200">
        <circle cx="10" cy="10" r="10"            fill="#456" stroke="none" />
        <circle cx="726" cy="180" r="10"          fill="#456" stroke="none" />
        <path   d="M 726 180 C 276 87 270 10 0 0" fill="none" stroke="#456"/>
    </svg>
    `  
  })
  export class ConnEndpointComponent implements OnInit {
    ngOnInit() {
        
      }

  }
