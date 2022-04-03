pragma solidity ^0.5.0;


contract ReferralBounties {

  /*
  * Enums
  */

  enum ReferralStatus { CREATED, ACCEPTED, REJECTED }

  /*
  * Storage
  */

  ReferralBounty[] public jobReferralList;

  mapping(uint=>ReferrerFulfillment[]) fulfillments;


  /*
  * Structs
  */

  struct ReferralBounty {
      address payable issuer;
      uint deadline;
      string data;
      ReferralStatus status;
      uint amount; 
  }

   address  payable constant reconnect = 0xE0ECb7C498F07D10a536bd9ac83CB5956cF70A8D ; 

  struct ReferrerFulfillment {
      bool accepted;
      address payable referrer;
      address  payable reconnect;
      string data;
  }
 

  /**
   * @dev Contructor
   */
  constructor() public {}

  /**
  * @dev issueBounty(): instantiates a new bounty
  * @param _deadline the unix timestamp after which fulfillments will no longer be accepted
  * @param _data the requirements of the bounty
  */
  function issueReferral(
      string calldata _data,
      uint64 _deadline
  )
      external payable hasValue() validateDeadline(_deadline) returns (uint){
      jobReferralList.push(ReferralBounty(msg.sender, _deadline, _data, ReferralStatus.CREATED, msg.value));
      emit jobReferralIssued(jobReferralList.length - 1,msg.sender, msg.value, _data);
      return (jobReferralList.length - 1);
  }

  /**
  * @dev fulfillBounty(): submit a fulfillment for the given bounty
  * @param _jobReferralId the index of the bounty to be fufilled
  * @param _data the ipfs hash which contains evidence of the fufillment
  */
  function referralFulfilled(uint _jobReferralId, string memory _data)
    public
    jobReferExists(_jobReferralId)
    notIssuer(_jobReferralId)
    hasStatus(_jobReferralId, ReferralStatus.CREATED)
    isBeforeDeadline(_jobReferralId)
  {
    fulfillments[_jobReferralId].push(ReferrerFulfillment(false, msg.sender, reconnect,_data));
    emit ReferralFulfilled(_jobReferralId, msg.sender, (fulfillments[_jobReferralId].length - 1),_data);
  }




  /**
  * @dev acceptFulfillment(): accept a given fulfillment
  * @param _jobReferralId the index of the bounty
  * @param _fulfillmentId the index of the fulfillment being accepted
  */
  function candidateAccepted(uint _jobReferralId, uint _fulfillmentId)
      public
      jobReferExists(_jobReferralId)
      fulfillmentExists(_jobReferralId,_fulfillmentId)
      onlyIssuer(_jobReferralId)
      hasStatus(_jobReferralId, ReferralStatus.CREATED)
      fulfillmentNotYetAccepted(_jobReferralId, _fulfillmentId)
  {
      fulfillments[_jobReferralId][_fulfillmentId].accepted = true;
      jobReferralList[_jobReferralId].status = ReferralStatus.ACCEPTED;
      fulfillments[_jobReferralId][_fulfillmentId].referrer.transfer(jobReferralList[_jobReferralId].amount/5);
      fulfillments[_jobReferralId][_fulfillmentId].reconnect.transfer(jobReferralList[_jobReferralId].amount/5);
      emit CandidateAccepted(_jobReferralId, jobReferralList[_jobReferralId].issuer, fulfillments[_jobReferralId][_fulfillmentId].referrer, _fulfillmentId, jobReferralList[_jobReferralId].amount,reconnect);
  }



  /** @dev cancelReferralBounty(): cancels the bounty and send the funds back to the issuer
  * @param _jobReferralId the index of the bounty
  */
  function cancelReferralBounty(uint _jobReferralId)
      public
      jobReferExists(_jobReferralId)
      onlyIssuer(_jobReferralId)
      hasStatus(_jobReferralId, ReferralStatus.CREATED)
  {
      jobReferralList[_jobReferralId].status = ReferralStatus.REJECTED;
      jobReferralList[_jobReferralId].issuer.transfer(jobReferralList[_jobReferralId].amount);
      emit CandidateRejected(_jobReferralId, msg.sender, jobReferralList[_jobReferralId].amount);
  }

  /**
  * Modifiers
  */

  modifier hasValue() {
      require(msg.value > 0);
      _;
  }

  modifier jobReferExists(uint _jobReferralId){
    require(_jobReferralId < jobReferralList.length);
    _;
  }

  modifier fulfillmentExists(uint _jobReferralId, uint _fulfillmentId){
    require(_fulfillmentId < fulfillments[_jobReferralId].length);
    _;
  }

  modifier hasStatus(uint _jobReferralId, ReferralStatus _desiredStatus) {
    require(jobReferralList[_jobReferralId].status == _desiredStatus);
    _;
  }

  modifier onlyIssuer(uint _jobReferralId) {
      require(msg.sender == jobReferralList[_jobReferralId].issuer);
      _;
  }

  modifier notIssuer(uint _jobReferralId) {
      require(msg.sender != jobReferralList[_jobReferralId].issuer);
      _;
  }

  modifier fulfillmentNotYetAccepted(uint _jobReferralId, uint _fulfillmentId) {
    require(fulfillments[_jobReferralId][_fulfillmentId].accepted == false);
    _;
  }

  modifier validateDeadline(uint _newDeadline) {
      require(_newDeadline > now);
      _;
  }

  modifier isBeforeDeadline(uint _jobReferralId) {
    require(now < jobReferralList[_jobReferralId].deadline);
    _;
  }

  /**
  * Events
  */
  event jobReferralIssued(uint jobReferralId, address issuer, uint amount, string data);
  event ReferralFulfilled(uint jobReferralId, address referrer, uint fulfillment_id, string data);
  event CandidateAccepted(uint jobReferralId, address issuer, address referrer, uint indexed fulfillment_id, uint amount, address reconnect);
  event CandidateRejected(uint indexed jobReferralId, address indexed issuer, uint amount);
}