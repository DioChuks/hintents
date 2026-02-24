import * as vscode from 'vscode';
import { Trace, TraceStep } from './erstClient';

export class TraceTreeDataProvider implements vscode.TreeDataProvider<TraceItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<TraceItem | undefined | null | void> = new vscode.EventEmitter<TraceItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<TraceItem | undefined | null | void> = this._onDidChangeTreeData.event;

    private currentTrace: Trace | undefined;

    constructor() { }

    refresh(trace: Trace): void {
        this.currentTrace = trace;
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: TraceItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: TraceItem): Thenable<TraceItem[]> {
        if (!this.currentTrace) {
            return Promise.resolve([]);
        }

        if (element) {
            return Promise.resolve([]);
        } else {
            const states = this.currentTrace.states;
            return Promise.resolve(
                states.map((step, idx) => new TraceItem(step, idx > 0 ? states[idx - 1] : undefined))
            );
        }
    }
}

export class TraceItem extends vscode.TreeItem {
    public isCrossContractBoundary: boolean;

    constructor(
        public readonly step: TraceStep,
        previousStep?: TraceStep
    ) {
        super(
            `${step.step}: ${step.operation}${step.function ? ` (${step.function})` : ''}`,
            vscode.TreeItemCollapsibleState.None
        );

        this.isCrossContractBoundary = isCrossContractTransition(previousStep, step);

        this.tooltip = `${this.label}`;
        this.description = step.error
            ? `Error: ${step.error}`
            : this.isCrossContractBoundary
                ? `[boundary] ${previousStep?.contract_id} -> ${step.contract_id}`
                : '';
        this.contextValue = this.isCrossContractBoundary ? 'traceStepBoundary' : 'traceStep';

        if (step.error) {
            this.iconPath = new vscode.ThemeIcon('error', new vscode.ThemeColor('errorForeground'));
        } else if (this.isCrossContractBoundary) {
            this.iconPath = new vscode.ThemeIcon('git-compare', new vscode.ThemeColor('editorWarning.foreground'));
        } else {
            this.iconPath = new vscode.ThemeIcon('pass', new vscode.ThemeColor('debugIcon.startForeground'));
        }
    }
}

// isCrossContractTransition returns true when two consecutive steps belong to different contracts.
function isCrossContractTransition(prev: TraceStep | undefined, current: TraceStep): boolean {
    if (!prev || !prev.contract_id || !current.contract_id) {
        return false;
    }
    return prev.contract_id !== current.contract_id;
}
