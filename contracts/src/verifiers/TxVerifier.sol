//
// Copyright 2017 Christian Reitwiessner
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// 2019 OKIMS
//      ported to solidity 0.6
//      fixed linter warnings
//      added requiere error messages
//
//
// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.8.0 <=0.9;
library Pairing {
    struct G1Point {
        uint X;
        uint Y;
    }
    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint[2] X;
        uint[2] Y;
    }
    /// @return the generator of G1
    function P1() internal pure returns (G1Point memory) {
        return G1Point(1, 2);
    }
    /// @return the generator of G2
    function P2() internal pure returns (G2Point memory) {
        // Original code point
        return G2Point(
            [11559732032986387107991004021392285783925812861821192530917403151452391805634,
             10857046999023057135944570762232829481370756359578518086990519993285655852781],
            [4082367875863433681332203403145435568316851327593401208105741076214120093531,
             8495653923123431417604973247489272438418190587263600148770280649306958101930]
        );

/*
        // Changed by Jordi point
        return G2Point(
            [10857046999023057135944570762232829481370756359578518086990519993285655852781,
             11559732032986387107991004021392285783925812861821192530917403151452391805634],
            [8495653923123431417604973247489272438418190587263600148770280649306958101930,
             4082367875863433681332203403145435568316851327593401208105741076214120093531]
        );
*/
    }
    /// @return r the negation of p, i.e. p.addition(p.negate()) should be zero.
    function negate(G1Point memory p) internal pure returns (G1Point memory r) {
        // The prime q in the base field F_q for G1
        uint q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
        if (p.X == 0 && p.Y == 0)
            return G1Point(0, 0);
        return G1Point(p.X, q - (p.Y % q));
    }
    /// @return r the sum of two points of G1
    function addition(G1Point memory p1, G1Point memory p2) internal view returns (G1Point memory r) {
        uint[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }
        require(success,"pairing-add-failed");
    }
    /// @return r the product of a point on G1 and a scalar, i.e.
    /// p == p.scalar_mul(1) and p.addition(p) == p.scalar_mul(2) for all points p.
    function scalar_mul(G1Point memory p, uint s) internal view returns (G1Point memory r) {
        uint[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }
        require (success,"pairing-mul-failed");
    }
    /// @return the result of computing the pairing check
    /// e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
    /// For example pairing([P1(), P1().negate()], [P2(), P2()]) should
    /// return true.
    function pairing(G1Point[] memory p1, G2Point[] memory p2) internal view returns (bool) {
        require(p1.length == p2.length,"pairing-lengths-failed");
        uint elements = p1.length;
        uint inputSize = elements * 6;
        uint[] memory input = new uint[](inputSize);
        for (uint i = 0; i < elements; i++)
        {
            input[i * 6 + 0] = p1[i].X;
            input[i * 6 + 1] = p1[i].Y;
            input[i * 6 + 2] = p2[i].X[0];
            input[i * 6 + 3] = p2[i].X[1];
            input[i * 6 + 4] = p2[i].Y[0];
            input[i * 6 + 5] = p2[i].Y[1];
        }
        uint[1] memory out;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 8, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }
        require(success,"pairing-opcode-failed");
        return out[0] != 0;
    }
    /// Convenience method for a pairing check for two pairs.
    function pairingProd2(G1Point memory a1, G2Point memory a2, G1Point memory b1, G2Point memory b2) internal view returns (bool) {
        G1Point[] memory p1 = new G1Point[](2);
        G2Point[] memory p2 = new G2Point[](2);
        p1[0] = a1;
        p1[1] = b1;
        p2[0] = a2;
        p2[1] = b2;
        return pairing(p1, p2);
    }
    /// Convenience method for a pairing check for three pairs.
    function pairingProd3(
            G1Point memory a1, G2Point memory a2,
            G1Point memory b1, G2Point memory b2,
            G1Point memory c1, G2Point memory c2
    ) internal view returns (bool) {
        G1Point[] memory p1 = new G1Point[](3);
        G2Point[] memory p2 = new G2Point[](3);
        p1[0] = a1;
        p1[1] = b1;
        p1[2] = c1;
        p2[0] = a2;
        p2[1] = b2;
        p2[2] = c2;
        return pairing(p1, p2);
    }
    /// Convenience method for a pairing check for four pairs.
    function pairingProd4(
            G1Point memory a1, G2Point memory a2,
            G1Point memory b1, G2Point memory b2,
            G1Point memory c1, G2Point memory c2,
            G1Point memory d1, G2Point memory d2
    ) internal view returns (bool) {
        G1Point[] memory p1 = new G1Point[](4);
        G2Point[] memory p2 = new G2Point[](4);
        p1[0] = a1;
        p1[1] = b1;
        p1[2] = c1;
        p1[3] = d1;
        p2[0] = a2;
        p2[1] = b2;
        p2[2] = c2;
        p2[3] = d2;
        return pairing(p1, p2);
    }
}
contract TxVerifier {
    using Pairing for *;
    struct VerifyingKey {
        Pairing.G1Point alfa1;
        Pairing.G2Point beta2;
        Pairing.G2Point gamma2;
        Pairing.G2Point delta2;
        Pairing.G1Point[] IC;
    }
    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }
    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            709611692107895393502317390195592908174896616524944250068206498856322549186,
            15894151335670580693397980132720711124065225275611250691534706720305437805427
        );

        vk.beta2 = Pairing.G2Point(
            [6784485304059443351313230882654010266656492193911369766268185595217537125805,
             8736214389600085643982110626358489434190961915445614650312313684171412410573],
            [7277654859118687732490630556653515568271209376238913180972984827580516425108,
             10114723301049440133092159159774945211340036895345626353973804265294610774028]
        );
        vk.gamma2 = Pairing.G2Point(
            [11559732032986387107991004021392285783925812861821192530917403151452391805634,
             10857046999023057135944570762232829481370756359578518086990519993285655852781],
            [4082367875863433681332203403145435568316851327593401208105741076214120093531,
             8495653923123431417604973247489272438418190587263600148770280649306958101930]
        );
        vk.delta2 = Pairing.G2Point(
            [11559732032986387107991004021392285783925812861821192530917403151452391805634,
             10857046999023057135944570762232829481370756359578518086990519993285655852781],
            [4082367875863433681332203403145435568316851327593401208105741076214120093531,
             8495653923123431417604973247489272438418190587263600148770280649306958101930]
        );
        vk.IC = new Pairing.G1Point[](20);
        
        vk.IC[0] = Pairing.G1Point( 
            20631919969101928388754347491964925772218565177656464713327144056617163403417,
            1404748013895498510343750067175631050109445248204240477469501170711797719326
        );                                      
        
        vk.IC[1] = Pairing.G1Point( 
            5618633452056373028991386307858651133956055669272002750789951885532750547508,
            7034991569263553625732010493486743595664666280647313711752396257289481331930
        );                                      
        
        vk.IC[2] = Pairing.G1Point( 
            1399765671306293060085935809849933042563875298042612598479191501251939146968,
            1534634028339201266829942320313128534270227172683431269524293939044448565646
        );                                      
        
        vk.IC[3] = Pairing.G1Point( 
            15338556109127290966798335390871286371904477075062876785889883985124614521086,
            4701458304165691185493595985561795321397133211536946344782765120315371276575
        );                                      
        
        vk.IC[4] = Pairing.G1Point( 
            1486699729880755631983687909733246342248719786538953763542514477882056505761,
            15010192011871871126787370643355387178423003743628273239235716041401883210369
        );                                      
        
        vk.IC[5] = Pairing.G1Point( 
            14802197714565958229214740053166636865064209824913211598148407751034848709571,
            17521539448394203748172537136129551276205535970458562734404591144782677643191
        );                                      
        
        vk.IC[6] = Pairing.G1Point( 
            20721824051198815887613393801568904167423285614504936692072399136536671350852,
            18831168153213426590617100109026534907620304447748991129024321472207529173184
        );                                      
        
        vk.IC[7] = Pairing.G1Point( 
            19828704235466902091619995369303822895897229832087814567270130306451612935433,
            2684314927436807221616476890112724886679395806857347823739536442126531959956
        );                                      
        
        vk.IC[8] = Pairing.G1Point( 
            9632595216027571175707579614744832218311073748502605294162913741034874802315,
            12479148059549787458866117404557901781066351312909143587512629465354151203758
        );                                      
        
        vk.IC[9] = Pairing.G1Point( 
            13206609993348650498636644041163370576274425285827361560017113801531487719757,
            11403675051390144401562077896139275437211900269379695838311104217402891119240
        );                                      
        
        vk.IC[10] = Pairing.G1Point( 
            11358096218454485202971108524483479794368356894660435036668165841754918518325,
            7800977582498906413124778888008549654774146137295466995527242177938476598104
        );                                      
        
        vk.IC[11] = Pairing.G1Point( 
            13471051542987120745384367993556134635190883921921181930639736571375155882604,
            19764355271328819153129994562241429707845552508891739469392610444249344848294
        );                                      
        
        vk.IC[12] = Pairing.G1Point( 
            7741590194671812015928894231572598898650118314573890989350908011626740410941,
            6097889402624584817684060373545204861350112724967939259348701595201960145928
        );                                      
        
        vk.IC[13] = Pairing.G1Point( 
            14533216280892226052606411538833381823538464363840546716643144785335735307659,
            15616498990980967797933080095770318611923338093526511343640392502953620236705
        );                                      
        
        vk.IC[14] = Pairing.G1Point( 
            4268092321459304505283265811481531439110900326836811038112048805498726644141,
            10453464944477963591974504060130845952701790470151642406692867722858926013493
        );                                      
        
        vk.IC[15] = Pairing.G1Point( 
            14116777456978590012269964550615371051618296875483412381615854319904223773367,
            3630796481274146057115734287032313032916179522745314858836893091958404335454
        );                                      
        
        vk.IC[16] = Pairing.G1Point( 
            4028600120233292032217339542160796955607425202677403825528498130261848249728,
            20113181009446062589503938436843635781671615719703868990606975400571185447567
        );                                      
        
        vk.IC[17] = Pairing.G1Point( 
            6019947519611459259304828089069239947361964487896617121203519899302006494990,
            12258107513142619390359305705188533358488630519948842650248044921104149668743
        );                                      
        
        vk.IC[18] = Pairing.G1Point( 
            15417281369967529998827649802477818010710033154942642823422416258087754115186,
            17925436801902230948983066889916470107621582937807158231692028668868275243672
        );                                      
        
        vk.IC[19] = Pairing.G1Point( 
            8640768230820827506367386294666025522270513261058139135512414662314278680239,
            7161788057697949612270884575602648914460576609498223956884344875516623433552
        );                                      
        
    }
    function verify(uint[] memory input, Proof memory proof) internal view returns (uint) {
        uint256 snark_scalar_field = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
        VerifyingKey memory vk = verifyingKey();
        require(input.length + 1 == vk.IC.length,"verifier-bad-input");
        // Compute the linear combination vk_x
        Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);
        for (uint i = 0; i < input.length; i++) {
            require(input[i] < snark_scalar_field,"verifier-gte-snark-scalar-field");
            vk_x = Pairing.addition(vk_x, Pairing.scalar_mul(vk.IC[i + 1], input[i]));
        }
        vk_x = Pairing.addition(vk_x, vk.IC[0]);
        if (!Pairing.pairingProd4(
            Pairing.negate(proof.A), proof.B,
            vk.alfa1, vk.beta2,
            vk_x, vk.gamma2,
            proof.C, vk.delta2
        )) return 1;
        return 0;
    }
    /// @return r  bool true if proof is valid
    function verifyProof(
            uint[2] memory a,
            uint[2][2] memory b,
            uint[2] memory c,
            uint[19] memory input
        ) public view returns (bool r) {
        Proof memory proof;
        proof.A = Pairing.G1Point(a[0], a[1]);
        proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
        proof.C = Pairing.G1Point(c[0], c[1]);
        uint[] memory inputValues = new uint[](input.length);
        for(uint i = 0; i < input.length; i++){
            inputValues[i] = input[i];
        }
        if (verify(inputValues, proof) == 0) {
            return true;
        } else {
            return false;
        }
    }
}
