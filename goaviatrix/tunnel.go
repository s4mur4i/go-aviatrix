package goaviatrix

// Tunnel simple struct to hold tunnel details

import (
	"fmt"
	"encoding/json"
	"errors"
	//"io/ioutil"
	"github.com/davecgh/go-spew/spew"
)

type Tunnel struct {
	VpcName1        string `json:"vpc_name1"`
	VpcName2        string `json:"vpc_name2"`
	OverAwsPeering  string `json:"over_aws_peering"`
	PeeringState    string `json:"peering_state"`
	PeeringHaStatus string `json:"peering_ha_status"`
	Cluster         string `json:"cluster"`
	PeeringLink     string `json:"peering_link"`
}

type TunnelResult struct {
	PairList        []Tunnel `json:"pair_list"`
}

type TunnelListResp struct {
	Return  bool   `json:"return"`
	Results TunnelResult `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateTunnel(tunnel *Tunnel) (error) {
	path := c.baseUrl + fmt.Sprintf("?CID=%s&action=peer_vpc_pair&vpc_name1=%s&vpc_name2=%s", c.CID, tunnel.VpcName1, tunnel.VpcName2)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data ApiResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetTunnel(tunnel *Tunnel) (*Tunnel, error) {
	//tunnel.CID=c.CID
	path := c.baseUrl + fmt.Sprintf("?CID=%s&action=list_peer_vpc_pairs", c.CID)
	fmt.Println("PaTh: ", path)
	resp,err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data TunnelListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	spew.Dump(data)
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	tunList:= data.Results.PairList
	for i := range tunList {
    	if tunList[i].VpcName1 == tunnel.VpcName1 && tunList[i].VpcName2 == tunnel.VpcName2 {
        	return &tunList[i], nil
    	}
	}
	return nil, errors.New(fmt.Sprintf("Tunnel with gateways %s and %s not found.", tunnel.VpcName1, tunnel.VpcName2))	
}

func (c *Client) UpdateTunnel(tunnel *Tunnel) (error) {
	return nil
}

func (c *Client) DeleteTunnel(tunnel *Tunnel) (error) {
	//tunnel.CID=c.CID
	path := c.baseUrl + fmt.Sprintf("?CID=%s&action=unpeer_vpc_pair&vpc_name1=%s&vpc_name2=%s", c.CID, tunnel.VpcName1, tunnel.VpcName2)
	fmt.Println("PaTh: ", path)
	resp,err := c.Delete(path, nil)

	if err != nil {
		return err
	}
	var data ApiResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	fmt.Println(data.Results)
	return nil
}