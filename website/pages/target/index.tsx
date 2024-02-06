'use client'
import React, { Component,useEffect, useRef ,useState} from "react";
import * as d3 from "d3";

function Target(){
  const ref = useRef<HTMLDivElement>(null);
  interface MySimulationNodeDatum extends d3.SimulationNodeDatum{
    name?: string;
    id:string;
    title:string;
    properties?:any;
    entity?:string;
  }
  interface MySimulationLinkDatum extends d3.SimulationLinkDatum<d3.SimulationNodeDatum>{
    relation?: string;
    value?:string;
    label?:string;
  }
  const [formData, setFormData] = useState({ field: '', keyword: '' ,rel:"",start:"",end:""});
  const handleChange = (e:any) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  async function getData(e:any){
    e.preventDefault();
    const res = await fetch("/api/target_graph")
    res.json().then((data)=>{
      // d3.select(ref.current).select("svg").remove()
      const nodes: MySimulationNodeDatum[] = [] ;
      
      const links: MySimulationLinkDatum[] = [];
      data.nodes.forEach((item:any)=>{
        nodes.push(item)
      })
      data.edges.forEach((item:any)=>{
        links.push({source:item.from,target:item.to,relation:item.label})
      })
      console.log(nodes)
      console.log(links)

      const  width=1000 
      const  height=800 
      
      const svg = d3.select(ref.current).select("svg")
        .attr("width",width)
        .attr("height",height)
        .attr("class","nodes")
        // .selectAll("*").remove()

      // 定义颜色比例尺
      const colorScale = d3.scaleOrdinal().domain(["User", "Product", "同学"]).range(["#815989", "#ebb9af", "#692b2b"]);

      // 创建力导图模拟
      const simulation = d3.forceSimulation().nodes(nodes)
        .force("link", d3.forceLink().links(links).id((d:any)=>d.id).distance(150))
        .force("charge", d3.forceManyBody().strength(-800))
        .force("collide", d3.forceCollide().radius(40))
        .force("center", d3.forceCenter(width / 2, height/ 2))
        .on("tick", () => {
          link.attr("x1",(d:any)=>{return d.source.x;})
            .attr("y1",(d:any)=>{return d.source.y;})
            .attr("x2",(d:any)=>{return d.target.x;})
            .attr("y2",(d:any)=>{return d.target.y;});

          linksText.attr("x",(d:any)=>{
            return (d.source.x+d.target.x)/2;
          }).attr("y",(d:any)=>{
            return (d.source.y+d.target.y)/2;
          });

          gs.attr("transform",function(d) { return "translate(" + d.x + "," + d.y + ")"; });
        });

      // 开始拖动
      function started(event:any, d:any) {
        if (!event.active) simulation.alphaTarget(0.1).restart();
        d.fx = d.x;
        d.fy = d.y;
      }

      // 拖动中
      function dragged(event:any, d:any) {
        d.fx = event.x;
        d.fy = event.y;
      }

      // 结束拖动
      function ended(event:any, d:any) {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
      }

      const g=svg.append("g").attr("class", "g");

      //绘制边
      const link = g.append("g")
        .selectAll("line")
        .data(links)
        .enter()
        .append("line")
        .attr("stroke",(d:any,i:any)=>{
          return "#eee";
        })
        .attr("stroke-width",2);

      const linksText = g.append("g")
          .selectAll("text")
          .data(links)
          .enter()
          .append("text")
          .text((d:any)=>{
            return d.relation;
          })
      const gs = g.selectAll(".circleText")
          .data(nodes)
          .enter()
          .append("g")
          .attr("transform",function(d,i){
            var cirX = d.x;
            var cirY = d.y;
            return "translate("+cirX+","+cirY+")";
          })
          .call(d3.drag<SVGGElement, MySimulationNodeDatum>()
            .on("start",started)
            .on("drag",dragged)
            .on("end",ended));
      //绘制节点
      gs.append("circle").attr("r",20).attr("refX",35).attr("refY",0).attr("fill",(d:any,i:any)=>{
        return colorScale(i) as string;
      })
      //文字
      gs.append("text")
        .attr("x",-14).attr("y",-4).attr("dy",8)
        .attr("class","text-xs text-white")
        .attr("stroke","#eee")
        .attr("stroke-width",1)
        .text((d:any)=>{
          return d.title;
        })
    })


  }

  useEffect(() => {
    

  },[]);

  return (
    <div>
      <div id="search" className="flex justify-start">
        <select id="field" name="country" onChange={handleChange} className="w-32 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block w-48 rounded-md sm:text-sm focus:ring-1">
          <option>资金账号</option>
          <option>手机号</option>
          <option>产品编号</option>
          <option>产品名称</option>
          <option>IP</option>
          <option>地址</option>
        </select>
        <input name="keyword" type="text" 
        className="ml-2 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block w-48 rounded-md sm:text-sm focus:ring-1" placeholder="请输入条件"/>
        <select id="rel" name="rel" onChange={handleChange} className="ml-2  px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block w-48 rounded-md sm:text-sm focus:ring-1">
          <option>推荐</option>
          <option>持有</option>
          <option>地址</option>
          <option>登录</option>
        </select>            
        <input name="start" type="text" onChange={handleChange} className="w-24 ml-2 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block rounded-md sm:text-sm focus:ring-1" placeholder="最小深度"/>
        <input name="end" type="text" onChange={handleChange} className="w-24 ml-2 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block rounded-md sm:text-sm focus:ring-1" placeholder="最小深度"/>
        <button type="button" onClick={getData} className="ml-2 px-3 py-2 justify-center rounded-md bg-red-500 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">搜索</button>
      </div>
      <div id="data" ref={ref} className="w-full py-4">
        <svg></svg>
      </div>
    </div>
  )
}

export default Target